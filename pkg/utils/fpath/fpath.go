package fpath

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func CurrentFileDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
