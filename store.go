package main

import (
	"io"
	"log"
	"os"
)

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

	fileName := "somefilename"
	pathAndFileName := pathName + "/" + fileName

	f, err := os.Create(pathAndFileName);
	if err != nil {
		return nil
	}
	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}

	log.Printf("dfs: wrote (%d) bytes to disk %s\n", n, pathAndFileName)

	return nil
}