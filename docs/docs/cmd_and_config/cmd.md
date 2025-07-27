## Cmd
This package is responsible for the interaction of console commands with the project.
It must be initialized by running the `Run` method.

### Initialization
For initial initialization, you need to create a file such as `cmd.go` with the following content:
```golang
package main

import "github.com/uwine4850/foozy/pkg/cmd"

func main() {
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
```
__Note__: most often, after import the `github.com/uwine4850/foozy/pkg/cmd` package, you need to run the `go mod tidy` command.

Commands can now be used with the following call: `go run cmd.go <command>`.

### Commands
Commands that are available in the package.

* cnf-init — initializes the configuration settings. The following file(including directory) is generated: `init_cnf/init_cnf.go`. The file has the following contents:
```golang
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
```
* cnf-gen — generates the `config.yaml` configuration file. To do this, the following operations must be performed:
    * Call the `cnf-init` command.
    * Call the `InitCnf()` function in the previously created `cmd.go` file. More information about the `InitCnf()` function is written [here](config.md#init-cnf). Now the `cmd.go` file should look like this:
```golang
package main

import (
    initcnf "tee/cnf/init_cnf"
    
    "github.com/uwine4850/foozy/pkg/cmd"
    )

func main() {
    initcnf.InitCnf()
    if err := cmd.Run(); err != nil {
        panic(err)
    }
}
```
* cnf-info — shows information about the configuration file.