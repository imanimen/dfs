package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

// CASPathTransformFunc takes a key string and returns a string representation of a file path
// based on a SHA1 hash of the key. The path is constructed by splitting the hash into
// 5-character segments and joining them with forward slashes.
// This function is used to generate a file path for storing data in a content-addressable
// storage (CAS) system, where the path is derived from the content of the data being stored.
func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}
	return strings.Join(paths, "/")
}

type PathTransformFunc func(string) string

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

func (s Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
	  return err
	}
  
	data, err := io.ReadAll(r)
	if err != nil {
	  return err
	}
  
	filenameBytes := md5.Sum(data)
	filename := hex.EncodeToString(filenameBytes[:])
	pathAndFileName := pathName + "/" + filename
  
	f, err := os.Create(pathAndFileName);
	if err != nil {
	  return err
	}
  
	n, err := io.Copy(f, bytes.NewReader(data)) 

	if err != nil {
	  return err
	}
  
	log.Printf("dfs: wrote (%d) bytes to disk %s\n", n, pathAndFileName)
  
	return nil
  }