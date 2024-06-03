package main

import (
	// "bytes"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	// "log"
	// "os"
	"strings"
)

// CASPathTransformFunc takes a key string and returns a string representation of a file path
// based on a SHA1 hash of the key. The path is constructed by splitting the hash into
// 5-character segments and joining them with forward slashes.
// This function is used to generate a file path for storing data in a content-addressable
// storage (CAS) system, where the path is derived from the content of the data being stored.
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s/", p.PathName, p.FileName)
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	_, err := os.Stat(pathKey.FullPath())
	return err == nil
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("dfs: deleted [%s]", pathKey.FileName)
	}()
	return os.RemoveAll(pathKey.FirstPathName())
}
func (s *Store) Read(key string) (io.Reader, error) {
	file, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	// Get the project root directory
	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	pathKey := s.PathTransformFunc(key)
	fullPath := filepath.Join(projectRoot, pathKey.FullPath())

	// Check if the directory structure already exists
	_, err = os.Stat(filepath.Dir(fullPath))
	if os.IsNotExist(err) {
		// Create the directory structure
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("dfs: wrote (%d) bytes to disk %s\n", n, fullPath)

	return nil
}