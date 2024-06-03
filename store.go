package main

import (
	// "bytes"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	// "path/filepath"

	// "log"
	// "os"
	"strings"
)

const DefaultRootFolderName = "dfs"

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
	// Root is the folder name of the root, containing all the files and folders.
	Root 			  string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = DefaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	_, err := os.Stat(pathKey.FullPath())
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())

	defer func() {
		log.Printf("dfs: deleted [%s]", fullPathWithRoot)
	}()

	return os.RemoveAll(fullPathWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	file, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	return buf, err
}

func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	FullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(FullPathWithRoot)
}


func (s *Store) WriteStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	file, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}
	n, err := io.Copy(file, r)
	if err != nil {
		return err
	}
	log.Printf("dfs: wrote (%d) bytes to disk %s\n", n, fullPathWithRoot)
	return nil
}