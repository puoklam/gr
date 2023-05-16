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
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fmt.Println(path)
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
	filename := strings.Replace(path, src, dst, -1)
	return os.WriteFile(filename, cl, 0666)
}

func Generate(dir string, dst string, files FileMap) error {
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return GenerateFile(path, dir, dst, files[path])
		}
		return nil
	}

	// convert slashes to os specific separator
	// remove trailing slashes
	dir = filepath.Clean(dir)
	dst = filepath.Clean(dst)
	return filepath.WalkDir(dir, fs.WalkDirFunc(walk))
}
