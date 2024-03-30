package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Return last part of path
// Trim right side path separator and split by separator
func LastPartOfPath(path string) string {
	sep := string(os.PathSeparator)
	path = strings.TrimRight(path, sep)
	components := strings.Split(path, sep)
	return components[len(components)-1]
}

// Get the list of files in the directory
func Dirwalk(dir string, depth int) []string {
	if depth < 0 {
		return []string{}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			subPaths := Dirwalk(filepath.Join(dir, file.Name()), depth-1)
			paths = append(paths, subPaths...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func CopyFile(src, dest string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(d, s)
	return err
}

func ExistsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func SameFile(path1, path2 string) bool {
	fileInfo1, err := os.Stat(path1)
	if err != nil {
		return false
	}
	fileInfo2, err := os.Stat(path2)
	if err != nil {
		return false
	}
	if os.SameFile(fileInfo1, fileInfo2) {
		return true
	}
	return false
}
