[UA](https://github.com/uwine4850/foozy/blob/master/docs/ua/ua_readme.md) | [EN](https://github.com/uwine4850/foozy)<br>
__foozy__ is a lightweight and flexible web framework. The project is based on the http.ServeMux and http.Server modules.
Also, modules depend on interfaces whenever possible, so most of them are open to change.

Modules that the framework contains: <br>
* builtin - built-in ready-made functionality, for example, authentication. It is not necessary to use.
* [database](https://github.com/uwine4850/foozy/blob/master/docs/en/database.md) - interface for working with the mysql database.
* interfaces - all golang interfaces used in the project.
* [livereload](https://github.com/uwine4850/foozy/blob/master/docs/en/livereload.md) - a module that can be used to restart the project after updating the files.
* [middlewares](https://github.com/uwine4850/foozy/blob/master/docs/en/middlewares.md) - module for creating middleware.
* [router](https://github.com/uwine4850/foozy/blob/master/docs/en/router.md) is the most important module, with the help of its functionality, project routing and much more are implemented.
* [form](https://github.com/uwine4850/foozy/blob/master/docs/en/form.md) - work with HTML forms.
* [server](https://github.com/uwine4850/foozy/blob/master/docs/en/server.md) - an add-on over http. Server for easier use and work with the router module.
* tmlengine - project templating engine. The pongo2 library is used.
* utils - general auxiliary functionality, for example, CSRF token generation.
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
The ``NewRouter(manager interfaces2.IManager) *Router`` method needs a manager to work, so the code will look like this:
```
newManager := router.NewManager()
newRouter := router.NewRouter(newManager)
```
The manager, in turn, needs the templating engine ``NewManager(engine interfaces2.ITemplateEngine) *Manager`` to work.
You need to add it:
```
newTmplEngine, err := tmlengine.NewTemplateEngine()
if err != nil {
    panic(err)
}
newManager := router.NewManager(newTmplEngine)
newRouter := router.NewRouter(newManager)
```
Next, you need to set the routes for the work. For example, to go to page __/home__ you need to make the following handler:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    manager.SetTemplatePath("templates/home.html")
    err := manager.RenderTemplate(w, r)
    if err != nil {
        panic(err)
    }
})
```
This code runs the handler at the address __/home __. When the user clicks on it, he will receive an HTML template by address
__ templates/home.html __, this template is set using "manager.SetTemplatePath (" templates/home.html ")," then
it is displayed using "manager.RenderTemplate (w, r)." <br>
It is important to note that it is not necessary to use a template (you can also change it), you can use
arranged method "w. Write ()" or others. For example, it is possible to display data on the page in JSON format, for example:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    values := map[string]string{"key1": "val1"}
    err = manager.RenderJson(values, w)
    if err != nil {
        panic(err)
    }
}
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
_server := server.NewServer(":8000", newRouter)
err = _server.Start()
if err != nil {
    panic(err)
}
```
This code means that the server will be running on the local host and will be on port 8000. To display pages
the router from the "newRouter" variable will be used. <br>
The full code of the mini-project is given below.
```
package main

import (
    "github.com/uwine4850/foozy/pkg/interfaces"
    "github.com/uwine4850/foozy/pkg/router"
    "github.com/uwine4850/foozy/pkg/server"
    "github.com/uwine4850/foozy/pkg/tmlengine"
    "net/http"
)

func main() {
    newTmplEngine, err := tmlengine.NewTemplateEngine()
    if err != nil {
        panic(err)
    }
    newManager := router.NewManager(newTmplEngine)
    newRouter := router.NewRouter(newManager)
    newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
        manager.SetTemplatePath("templates/home.html")
        err := manager.RenderTemplate(w, r)
        if err != nil {
            panic(err)
        }
    })
    newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    
    _server := server.NewServer(":8000", newRouter)
    err = _server.Start()
    if err != nil {
        panic(err)
    }
}
```