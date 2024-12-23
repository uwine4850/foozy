package initcnf

import (
	"github.com/uwine4850/foozy/pkg/config"
)

// InitCnf initializes the configuration settings.
// This function should be used before initializing cmd(Run() function)
// and the server.
func InitCnf() {
	cnf := config.Cnf()
	cnf.SetPath("config.yaml")
	cnf.SetLoadPath("../config.yaml")
}
