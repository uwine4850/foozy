package fpath

import (
	"os"
	"path/filepath"
	"runtime"
)

// CurrentFileDir returns the path to the directory in which this function is called.
func CurrentFileDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

// PathExist checks to see if a path exists in the file directory.
func PathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
