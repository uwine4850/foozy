package initcnf_t

import (
	"os"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/config"
)

// InitCnf initializes the configuration settings.
// This function should be used before initializing cmd(Run() function)
// and the server.
func InitCnf() {
	projectRoot := os.Getenv("FOOZY_PROJECT_ROOT")
	cnf := config.Cnf()
	// Use this to add your configurations.
	// cnf.AppendAdditionally("my_cnf", typeopr.Ptr{}.New(&MyCnfCtruct{}))
	cnf.SetPath(filepath.Join(projectRoot, "tests1/cnf/config.yaml"))
	cnf.SetLoadPath(filepath.Join(projectRoot, "tests1/common/cnf/config.yaml"))
}
