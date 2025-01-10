package initcnf

import (
	"github.com/uwine4850/foozy/pkg/config"
)

// InitCnf initializes the configuration settings.
// This function should be used before initializing cmd(Run() function)
// and the server.
func InitCnf() {
	cnf := config.Cnf()
	// Use this to add your configurations.
	// cnf.AppendAdditionally("my_cnf", typeopr.Ptr{}.New(&MyCnfCtruct{}))
	cnf.SetPath("cnf/config.yaml")
	cnf.SetLoadPath("cnf/config.yaml")
}
