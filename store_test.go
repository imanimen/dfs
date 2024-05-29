package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpricture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalName := "35e11881d4c98f52c39b5e17e98d0322beda8feb78a03b0d8181da4b6d9c208fdd542c81"
	exptPathName := "35e11/881d4/c98f5/2c39b/5e17e/98d03/22bed/a8feb/78a03b0d8181da4b6d9c208fdd542c81"
	if pathKey.PathName != exptPathName {
		t.Errorf("dfs: expected pathName to be %s, but got %s", pathKey.PathName, exptPathName)
	}

	if pathKey.PathName != exptPathName {
		t.Errorf("dfs: expected pathName to be %s, but got %s", pathKey.PathName, expectedOriginalName)
	}

}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg byte"))
	if err := s.writeStream("mySpecialPicture", data); err != nil {
		t.Error(err)
	}
	
}