package fpath

import (
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
