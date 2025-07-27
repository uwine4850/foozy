## Getting started
This section will describe in detail how to get started with the framework and create a minimally working project. 

You can read more about all the packages used in the rest of the documentation.

## Install
You need to use the appropriate command to install the framework:
```bash
go get github.com/uwine4850/foozy
```
Naturally you need to have `golang 1.20` or more installed.

Also, after using some packages you need to install dependencies. To do this, just use the `go mod tidy` command.

## Creating a project
The following will describe the successive steps in creating a project. They must be performed in sequence to minimally initialize the project.

### Commands and configuration
To create a project you need to use several console commands. To do this you need to initialize the [cmd package](cmd_and_config/cmd.md).

One way will be shown below, it is just one implementation, do not consider it a single correct implementation.

* Create a `command` directory.
* Create a `commad.go` file.
* Populate the `command.go` file with the following code:
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

* Call the `go run command/command.go cnf-init ./cnf` [command](cmd_and_config/cmd.md#commands). Generating the `cnf/cnf_init/cnf_init.go` file.
* From the generated file, you need to import the [initcnf.InitCnf() function](cmd_and_config/config.md#init-cnf). The `command/command.go` file now looks like this:
```golang
package main

import (
	initcnf "<project-mod>/cnf/init_cnf"

	"github.com/uwine4850/foozy/pkg/cmd"
)

func main() {
	initcnf.InitCnf()
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
```
This is to know the path to the configuration file.

* Next, you need to exec the `go run command/command.go cnf-gen` [command](cmd_and_config/cmd.md#commands). This will generate a configuration file here (if the user has not changed anything) `cnf/config.yaml`.

__End__

After these operations, the console commands and the cofiguration file are ready for use.
You can read more details here:

* [Console commands](cmd_and_config/cmd.md).
* [Configuration file](cmd_and_config/config.md).

### Launching the first page
To start, a complete example of a minimal server just for copying and a quick test will be shown. After startup you can see the page at this address `http://localhost:7000/page`. Start the server with the `go run <filename>.go` command.

__Note:__ Most likely you will need to call the `go mod tidy` command.

```golang
package main

import (
	"net/http"
	initcnf "tee/cnf/init_cnf"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/server"
)

func main() {
	initcnf.InitCnf()
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		nil,
		nil,
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/page",
		func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
			w.Write([]byte("My first page!"))
			return nil
		})
	newServer := server.NewServer(":7000", newRouter, nil)
	if err := newServer.Start(); err != nil {
		panic(err)
	}
}
```
___
Below is a more detailed breakdown of all of this.

Initializing the configuration for use in the project. This only needs to be done once, as the configuration just needs to be loaded. This means that you only need to call this function once for a single session.  It does not matter where this function is called, but it must be called for the session. Read more [here](cmd_and_config/config.md#init-cnf).
```golang
initcnf.InitCnf()
```
___

[Manager](/router/manager/manager/) initialization. This object is used in many places in the framework, but you need to initialize it here. It is mandatory to pass [manager.NewOneTimeData()](/router/manager/manager/#onetimedata). __TODO: add link__ [Render]() and __TODO: add link__ [DatabasePool]() are optional if you don't plan to use a templating engine or database.
```golang
newManager := manager.NewManager(
    manager.NewOneTimeData(),
	nil,
	nil,
)
```
But for future use, full initialization of the [manager](/router/manager/manager/) is recommended.
```golang
newRender, err := tmlengine.NewRender()
if err != nil {
    panic(err)
}
newManager := manager.NewManager(
	manager.NewOneTimeData(),
	newRender,
	database.NewDatabasePool(),
)
```
___
Initializing the middleware. Read more [here](/router/middlewares/middlewares).
```golang
newMiddlewares := middlewares.NewMiddlewares()
```
___
Initializing the [adapter](/router/router/#adapter). It is needed for starting and preliminary preparation of handlers.
```golang
newAdapter := router.NewAdapter(newManager, newMiddlewares)
```
___
Initialize the [router](/router/router/#router) to handle http routes.
```golang
newRouter := router.NewRouter(newAdapter)
```
___
A [handler](/router/router/) that will process the selected routes. 
An error is returned, which will be handled by a special method specified in the [adapter](/router/router/#adapter).
```golang
newRouter.Register(router.MethodGET, "/page",
	func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("My first page!"))
		return nil
	})
```
This shows the simplest possible output of data to a page using `w.Write([]byte(“My first page!”))`.
For more advanced output, you need to use a __TODO: add link__ [templating engine](). But before that you need to create HTML file and make __full__ initialization of the [manager](/router/manager/manager/), which is shown above.
```golang
newRouter.Register(router.MethodGET, "/page",
	func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		manager.Render().SetTemplatePath("index.html")
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
```
___
Starting a __TODO: add link__ [server]() to process http requests. Here processing is started without __TODO: add link__ [cors]().
```golang
newServer := server.NewServer(":7000", newRouter, nil)
if err := newServer.Start(); err != nil {
	panic(err)
}
```
___

This is the complete minimal initialization of the framework. 

For more detailed information you need to look at other documentation or __TODO: add link__ [example project]().