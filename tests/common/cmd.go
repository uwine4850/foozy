package main

import (
	"github.com/uwine4850/foozy/pkg/cmd"
	testinitcnf "github.com/uwine4850/foozy/tests/common/test_init_cnf"
)

func main() {
	testinitcnf.InitCnf()
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
