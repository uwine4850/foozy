package codegen

import (
	"io"
	"os"
	"path/filepath"
)

// Generate generates files from a defined path.
// The map key — path to the directory where the file will be.
// Map value — path to the file to be generated.
func Generate(data map[string]string) error {
	for dirpath, targetFilepath := range data {
		targetFile, err := os.Open(targetFilepath)
		if err != nil {
			return err
		}
		defer targetFile.Close()
		fullNewPath := filepath.Join(dirpath, filepath.Base(targetFilepath))
		if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
			return err
		}
		newFile, err := os.Create(fullNewPath)
		if err != nil {
			return err
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, targetFile)
		if err != nil {
			return err
		}
	}
	return nil
}
