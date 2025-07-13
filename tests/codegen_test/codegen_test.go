package codegen

import (
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/codegen"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

func TestGenerate(t *testing.T) {
	if err := codegen.Generate(map[string]string{
		"dir/": "file.txt",
	}); err != nil {
		t.Error(err)
	}
	fileOk, err := tutils.FilesAreEqual("file.txt", "dir/file.txt")
	if err != nil {
		t.Error(err)
	}
	if !fileOk {
		t.Error("filed don't match")
	}
	if err := os.RemoveAll("dir/"); err != nil {
		t.Error(err)
	}
}

func TestGenerateNoGenFile(t *testing.T) {
	if err := codegen.Generate(map[string]string{
		"dir/": "file1.txt",
	}); err != nil {
		if !os.IsNotExist(err) {
			t.Error(err)
		}
	}
}
