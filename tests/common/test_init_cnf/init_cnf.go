package testinitcnf

import (
	"os"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/config"
)

func InitCnf() {
	projectRoot := os.Getenv("FOOZY_PROJECT_ROOT")
	cnf := config.Cnf()
	cnf.SetPath(filepath.Join(projectRoot, "tests/common/config.yaml"))
	cnf.SetLoadPath(filepath.Join(projectRoot, "tests/common/config.yaml"))
}
