package cmd

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/codegen"
	"github.com/uwine4850/foozy/pkg/config"
)

var myArgs = map[string]func(args ...string) error{
	"cnf-info": cnfInfo,
	"cnf-init": cnfInit,
	"cnf-gen":  cnfGen,
}

// cnfInfo shows information about configuration fields.
func cnfInfo(args ...string) error {
	config.Info()
	return nil
}

// cnfInit initialization of config generation settings.
// Generates a file with configuration settings,
// only in it you need to change the configuration settings.
// cnf-init <target directory>
func cnfInit(args ...string) error {
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
}

// cnfGen generates a .yaml configuration file.
// It is important to note that it is advisable
// to initialize the configuration with the "cnf-init" command.
func cnfGen(args ...string) error {
	gen := config.NewGenerate(config.Cnf())
	if err := gen.Gen(); err != nil {
		return err
	}
	return nil
}

// Run runs cmd.
// For proper implementation, this function should be placed in the main package.
// Also after the “initcnf” command you should use “initcnf.InitCnf()”.
// Example of implementation:
//
//	func main() {
//		initcnf.InitCnf()
//		if err := cmd.Run(); err != nil {
//			panic(err)
//		}
//	}
func Run() error {
	flag.Parse()
	args := flag.Args()
	if err := myArgs[args[0]](args...); err != nil {
		return err
	}
	return nil
}
