package initcnf

import (
	"github.com/uwine4850/foozy/pkg/config"
)

func InitCnf() {
	cnf := config.Cnf()
	cnf.SetPath("config.yaml")
	cnf.SetLoadPath("../config.yaml")
}
