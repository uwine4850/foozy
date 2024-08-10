[UA](https://github.com/uwine4850/foozy/blob/master/docs/ua/ua_readme.md) | [EN](https://github.com/uwine4850/foozy)<br>
__foozy__ is a lightweight and flexible web framework. The project is based on the http.ServeMux and http.Server modules.
Also, modules depend on interfaces whenever possible, so most of them are open to change.

Example project - https://github.com/uwine4850/foozy_proj

Modules that the framework contains: <br>
* [builtin](https://github.com/uwine4850/foozy/blob/master/docs/en/builtin/builtin.md) — built-in ready-made functionality, for example, authentication. It is not necessary to use.
* [database](https://github.com/uwine4850/foozy/blob/master/docs/en/database/database.md) — an interface for working with the mysql database.
  * [dbutils](https://github.com/uwine4850/foozy/blob/master/docs/en/database/dbutils/dbutils.md) — auxiliary functionality for using the database package.
  * [dbmapper](https://github.com/uwine4850/foozy/blob/master/docs/en/database/dbmapper/dbmapper.md) — writes data to the selected object.
  * [sync_queres](https://github.com/uwine4850/foozy/blob/master/docs/en/database/sync_queries.md) — synchronous requests to the database.
  * [async_queres](https://github.com/uwine4850/foozy/blob/master/docs/en/database/async_queries.md) — asynchronous requests to the database.
* [interfaces](https://github.com/uwine4850/foozy/blob/master/docs/en/interfaces/interfaces.md) — all golang interfaces used in the project.
* [router](https://github.com/uwine4850/foozy/blob/master/docs/en/router/router.md) — is the most important module, with the help of its functionality project routing and much more is implemented.
  * [manager](https://github.com/uwine4850/foozy/blob/master/docs/en/router/manager/manager.md) — a package for managing processors.
  * [websocket](https://github.com/uwine4850/foozy/blob/master/docs/en/router/websocket.md) — a package for interaction with websockets.
  * [form](https://github.com/uwine4850/foozy/blob/master/docs/en/router/form/form.md) — working with HTML forms.
	* [formmapper](https://github.com/uwine4850/foozy/blob/master/docs/en/router/form/formmapper/formmapper.md) — various manipulations with form data.
  * [middlewares](https://github.com/uwine4850/foozy/blob/master/docs/en/router/middlewares/middlewares.md) — a module for creating middleware.
  * [object](https://github.com/uwine4850/foozy/blob/master/docs/en/router/object/object.md) — a package for simpler display of templates.
  * [mic](https://github.com/uwine4850/foozy/blob/master/docs/en/router/mic/mic.md) — package is responsible for the functionality of microservices.
  * [tmlengine](https://github.com/uwine4850/foozy/blob/master/docs/en/router/tmlengine/tmlengine.md) — project templater. The pongo2 library is used.
* [server](https://github.com/uwine4850/foozy/blob/master/docs/en/server/server.md) — add-on to http.Server for easier use and work with the router module.
  * [livereload](https://github.com/uwine4850/foozy/blob/master/docs/en/server/livereload/livereload.md) — a module that can be used to restart the project after updating the files.
* [utils](https://github.com/uwine4850/foozy/blob/master/docs/en/utils/utils.md) — general auxiliary functionality, for example, CSRF token generation.

## Getting started

### Installation
```
go get github.com/uwine4850/foozy
```

### Basic usage
First you need to use a router ``router.Router``, for example:
```
newRouter := router.NewRouter()
```
The ``NewRouter(manager interfaces.IManager) *Router`` method needs a manager to work, so the code will look like this:
```
newManager := router.NewManager()
newRouter := router.NewRouter(newManager)
```
In turn, the manager needs a render structure to work ``NewManager(render interfaces.IRender) *Manager`` to work.
You need to add it:
```
render, err := tmlengine.NewRender()
if err != nil {
    panic(err)
}
newManager := router.NewManager(render)
newRouter := router.NewRouter(newManager)
```
Next, you need to set the routes for the work. For example, to go to page __/home__ you need to make the following handler:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    manager.Render().SetTemplatePath("templates/home.html")
	if err := manager.Render().RenderTemplate(w, r); err != nil {
	    panic(err)
    }
    return func() {}
})
```
This code runs the handler at the address __/home __. When the user clicks on it, he will receive an HTML template by address
__ templates/home.html __, this template is set using ``manager.Render().SetTemplatePath("templates/home.html")``, then it is displayed using "manager.RenderTemplate (w, r)." <br>
It is important to note that it is not necessary to use a template (you can also change it), you can use
arranged method "w. Write ()" or others. For example, it is possible to display data on the page in JSON format, for example:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    values := map[string]string{"key1": "val1"}
	if err := manager.Render().RenderJson(values, w); err != nil {
		panic(err)
	}
	return func() {}
})
```
Every website requires CSS and JavaScript in addition to HTML. You can add it with the following code:
```
newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
```
This code assumes that all files will be located in the static directory. In HTML, the file connection will look like this:
```
<link rel="stylesheet" href="/static/css/style.css">
...
<img src="/static/img/image.png">
```
It is important to note that the path must always start with the character ``/``.<br>
So, now that you have the basic handler, you need to start the server as follows.
```
serv := server.NewServer(":8000", newRouter)
err = serv.Start()
if err != nil && !errors.Is(http.ErrServerClosed, err) {
	panic(err)
}
```
This code means that the server will be running on the local host and will be on port 8000. To display pages
the router from the "newRouter" variable will be used. <br>
The full code of the mini-project is given below.
```
package main

import (
	"errors"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
)

func main() {
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newManager := manager.NewManager(render)
	newRouter := router.NewRouter(newManager)
    newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) 
    func() {
        manager.Render().SetTemplatePath("templates/home.html")
        if err := manager.Render().RenderTemplate(w, r); err != nil {
            panic(err)
        }
        return func() {}
    })
	serv := server.NewServer(":8000", newRouter)
	err = serv.Start()
	if err != nil && !errors.Is(http.ErrServerClosed, err) {
		panic(err)
	}
}
```
