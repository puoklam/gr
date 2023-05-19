package files

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GenerateFile(path string, src string, dst string, vars []*Var) error {
	fmt.Println("Reading file:", path)
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	/*
		TODO: replace all vars at same time to prevent override
		e.g. {{$a}} replace with {{$b}}
			 {{$b}} replace with 1
			 then {{$a}} will be replaced with 1 whch is unexpected
	*/

	cl := bytes.Clone(bs)
	for _, v := range vars {
		cl = bytes.ReplaceAll(cl, []byte(v.Name), []byte(v.Temp))
	}
	for _, v := range vars {
		cl = bytes.ReplaceAll(cl, []byte(v.Temp), []byte(v.Replace))
	}
	var filename string
	if path == src {
		// src is a file
		_, file := filepath.Split(path)
		filename = filepath.Join(dst, file)
	} else {
		// src is a directory
		filename = strings.Replace(path, src, dst, -1)
	}
	fmt.Println("Writing file:", filename)
	return os.WriteFile(filename, cl, 0666)
}

func Generate(src string, dst string, files FileMap) error {
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return GenerateFile(path, src, dst, files[path])
		}
		return nil
	}

	// convert slashes to os specific separator
	// remove trailing slashes
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	return filepath.WalkDir(src, fs.WalkDirFunc(walk))
}
