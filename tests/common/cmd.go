package main

import (
	"github.com/uwine4850/foozy/pkg/cmd"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
)

func main() {
	initcnf_t.InitCnf()
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
