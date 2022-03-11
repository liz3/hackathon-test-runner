package main

import (
"path"
"io/ioutil"
)

type File struct {
	Name string
	Path string
	Content string
}

func ReadFile(target string) *File {
	content, err := ioutil.ReadFile(target)
	if err != nil {
		return nil
	}
	return &File{
		Name: path.Base(target),
		Path: target,
		Content: string(content),
	}
}
func WriteFile(target string, content []byte) {
    err := ioutil.WriteFile(target, content, 0644)
    if err != nil {
    	panic(err)
    }
}