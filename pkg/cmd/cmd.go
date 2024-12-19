package cmd

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/codegen"
	"github.com/uwine4850/foozy/pkg/config"
)

var myArgs = map[string]func(args ...string) error{
	"initcnf": func(args ...string) error {
		if len(args) != 2 {
			return errors.New("parent directory not specified")
		}
		genfiles := map[string]string{
			filepath.Join(args[1], "init_cnf"): "internal/codegen/init_cnf/init_cnf.go",
		}
		if err := codegen.Generate(genfiles); err != nil {
			return err
		}
		return nil
	},
	"gencnf": func(args ...string) error {
		gen := config.NewGenerate(config.Cnf())
		if err := gen.Gen(); err != nil {
			return err
		}
		return nil
	},
}

func Run() error {
	flag.Parse()
	args := flag.Args()
	if err := myArgs[args[0]](args...); err != nil {
		return err
	}
	return nil
}
