## livereload
This package implements a system for reloading after saving a file. It does not interact directly with the framework and is an optional component, but it makes development much easier.

Example of use:
```golang
package main

import (
	initcnf "github.com/uwine4850/foozy/mycmd/init_cnf"
	"github.com/uwine4850/foozy/pkg/server/livereload"
)

func main() {
	initcnf.InitCnf()
	wrt := livereload.NewWiretap()
	wrt.SetDirs([]string{"cnf"})
	reload := livereload.NewReloader("main.go", wrt)
	if err := reload.Start(); err != nil {
		panic(err)
	}
}
```

More details about the package components:

* [wiretap](/server/livereload/wiretap)
* [reload](/server/livereload/reload)