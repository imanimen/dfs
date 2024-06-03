package main

import (
	"bytes"
	"fmt"
	"io"

	// "fmt"
	// "io"
	"testing"
)

// func TestPathTransformFunc(t *testing.T) {
// 	key := "momsspecials"
// 	pathKey := CASPathTransformFunc(key)
// 	expectedOriginalKey := "ff254eed1e1731bb8327808fd4700135c58a3e91"
// 	expectedPathName := "ff254/eed1e/1731b/b8327/808fd4700135c58a3e91"
// 	if pathKey.PathName != expectedPathName {
// 		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
// 	}
// 	if pathKey.FileName != expectedPathName {
// 		t.Errorf("have %s want %s", pathKey.FileName, expectedOriginalKey)
// 	}
// }

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momsspecials"
	data := []byte("some jpg bytes")
	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	// fmt.Println(k)

	err := s.Delete(key)
	if err != nil {
		return
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momsspecials"
	data := []byte("some jpg bytes")
	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

}
