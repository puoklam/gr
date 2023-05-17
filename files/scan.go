package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type Var struct {
	Name    string   // name of the variable
	Temp    string   // temp replacement to solve the issue of overriding variables
	Replace string   // string to replce the variable
	Paths   []string // files that contain the variable
}

type FileMap = map[string][]*Var

func ScanDir(dir string) ([]*Var, FileMap, error) {
	// varname -> file path
	vm := make(map[string][]string)

	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		names, err := ScanFile(path)
		if err != nil {
			return err
		}
		for _, nb := range names {
			name := string(nb)
			vm[name] = append(vm[name], string(path))
		}
		return nil
	}

	err := filepath.WalkDir(dir, fs.WalkDirFunc(walk))
	if err != nil {
		return nil, nil, err
	}

	vars := make([]*Var, 0, len(vm))
	files := make(FileMap)
	for name, paths := range vm {
		v := &Var{
			Name: name,
		}
		for _, path := range paths {
			files[path] = append(files[path], v)
		}
		v.Paths = make([]string, len(paths))
		copy(v.Paths, paths)
		vars = append(vars, v)
	}

	// sort vars
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars, files, nil
}

func ScanFile(path string) ([][]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`\{\{\$[A-Za-z]+\}\}`)
	return re.FindAll(b, -1), nil
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
