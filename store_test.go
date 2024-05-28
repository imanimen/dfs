package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpricture"
	pathName := CASPathTransformFunc(key)
	fmt.Println(pathName)
	exptPathName := "8d28b/f64ba/39cf9/27180/490ed/86827/d71ae/b3526"
	if pathName != exptPathName {
		t.Errorf("dfs: expected pathName to be %s, but got %s", exptPathName, pathName)
	}

}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("mySpecialPicture", data); err != nil {
		t.Error(err)
	}
	
}