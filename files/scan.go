package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type VarFile struct {
	ref    *Var
	Path   string  // file path
	Occurs [][]int // occurences of the target var
}

type Var struct {
	Name    string    // name of the variable
	Replace string    // string to replce the variable
	Files   []VarFile // files that contain the variable
}

type FileMap = map[string][]*Var

func ScanDir(dir string) ([]*Var, FileMap, error) {
	varmap := make(map[string]map[string][][]int)
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			filemap, err := ScanFile(path)
			if err != nil {
				return err
			}
			for name, s := range filemap {
				if varmap[name] == nil {
					varmap[name] = make(map[string][][]int)
				}
				varmap[name][path] = s
			}
		}
		return nil
	}

	err := filepath.WalkDir(dir, fs.WalkDirFunc(walk))
	if err != nil {
		return nil, nil, err
	}

	vars := make([]*Var, 0, len(varmap))
	files := make(FileMap)
	for name, filemap := range varmap {
		varFiles := make([]VarFile, 0, len(filemap))
		ref := &Var{
			Name: name,
		}
		for path, s := range filemap {
			varFiles = append(varFiles, VarFile{
				ref:    ref,
				Path:   path,
				Occurs: s,
			})
			files[path] = append(files[path], ref)
		}
		ref.Files = varFiles
		vars = append(vars, ref)
	}
	return vars, files, nil
}

func ScanFile(path string) (map[string][][]int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`\{\{\$[A-Za-z]+\}\}`)
	return matchVars(b, re), nil
}

func matchVars(src []byte, re *regexp.Regexp) map[string][][]int {
	m := make(map[string][][]int)
	matches := re.FindAllIndex(src, -1)
	for _, idx := range matches {
		v := string(src[idx[0]:idx[1]])
		m[v] = append(m[v], idx)
	}
	return m
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
