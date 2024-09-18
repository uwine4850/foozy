## Router package

## Router
The router is responsible for routing and handling processors. Below are the methods that are available for use.<br>

Tests for the router [here](https://github.com/uwine4850/foozy/tree/master/tests/routing).

### Route handlers
All handlers have standard parameters:
* *pattern* - the path to the handler. For example, the path can be "/home" and everyone is like it. The router also
    supports slug parameters, such as "/post/<id>". If you use such a path and go to the address "/post/1" - it will start
    the required handler in which the __id__ option will be available in the manager like this ``manager.GetSlugParams ("id")``. Such slug
    parameters can be many, the main thing with different names.
* *fn func(w http.ResponseWriter, r \*http.Request, manager interfaces.IManager)* - a function that will run when the user
  goes to the desired address. ``w http.ResponseWriter`` and ``*http.Request`` are standard golang structures. About ``interfaces.IManager``
  is described in more detail [here](https://github.com/uwine4850/foozy/blob/master/docs/en/router/manager/manager.md).
* Each handler returns ``func()`` - this is a function that is executed after the handler itself is finished.

Multiple method handlers can be applied to a single route. 
The main thing is that the methods are not repeated, that is, only one Get, Post, Delete, etc.

__Get__
```
Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
The method is used to transfer data from the server to the web page. For example, it can be html data, JSON data and others.

__Post__
```
Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
A method for handling POST requests. It is most often used to work with HTML forms. This method has no new functionality 
compared to the ``Get`` method. But this can be changed with [form handler package](https://github.com/uwine4850/foozy/blob/master/docs/en/router/form/form.md).

__Ws__
```
Ws(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
The handler accepts a websocket via the selected path.  Web sockets are written in more detail [here](https://github.com/uwine4850/foozy/blob/master/docs/en/router/websocket.md).<br>
An example of echo handler implementation:
```
newRouter.Ws("/ws", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	ws := router.NewWebsocket(router.Upgrader)
	ws.OnConnect(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		fmt.Println("Connect.")
	})
	ws.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn websocket.Conn) {
		err := ws.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Client close.")
	})
	ws.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		err := ws.SendMessage(messageType, msgData, conn)
		if err != nil {
			panic(err)
		}
	})
	err = ws.ReceiveMessages(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
})
```

__Put__
```
Put(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func())
```
An example of echo handler implementation:

__Delete__
```
Delete(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
```
Method for handling Delete requests.

__Options__
```
Options(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
```
A method for processing Options requests.

__InternalError__
```
InternalError(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error))
```
InternalError sets the function to be used when handling internal errors.

### Other methods

__RegisterAll__
```
RegisterAll()
```
Registers all handlers.

__SetTemplateEngine__
```
SetTemplateEngine(engine interfaces.ITemplateEngine)
```
Installs a templater instance that implements the ``interfaces.ITemplateEngine`` interface.

__SetMiddleware__
```
SetMiddleware(middleware interfaces.IMiddleware)
```
Sets an instance of the interface ``interfaces.IMiddleware`` to run before each handler.

__getHandleFunc__
```
getHandleFunc(pattern string, method string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) http.HandlerFunc
```
A private method that is launched before each handler. It runs various validations, middleware and more.
This method wraps the ``func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)`` function passed from 
handler Post, Get or Ws.

__validateMethod__
```
validateMethod(method string) bool
```
Validation of http methods currently in use. That is, it makes sure that the Get method processes the http GET method, and not, for example, POST.

#### Methods that do not belong to an interface but belong to a package
These methods are global, found in the router package, but can be used anywhere.<br>

__ValidateRootUrl__
```
ValidateRootUrl(w http.ResponseWriter, r *http.Request) bool
```
If the path pattern is __/__, then the handler will accept __all paths__. To prevent this you need to use this 
method. Now if the path is not found, a 404 error will be displayed instead of the __/__ handler. Example:
```
newRouter.Get("/", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    if !router.ValidateRootUrl(w, r) {
	    return
	}
})
```
__ParseSlugIndex__
```
ParseSlugIndex(path []string) map[int]bool
```
The function divides the url into parts based on the __/__ character. Next, he adds each part to the map in numerical order as a key, and a value 
the key depends on whether it is a slug parameter, if so the key is true, if not - false.

__HandleSlugUrls__
```
HandleSlugUrls(parseUrl map[int]bool, slugUrl []string, url []string) (string, map[string]string)
```
``slugUrl []string`` this is a pattern parameter that is divided by the __/__ symbol and written in slice.
``url []string`` this is the real url that is divided by the __/__ symbol and written in a slice.

The function processes the url and outputs a string (str) as a url made from the pattern and slug parameters, if any.<br>
Use ``parseUrl map[int]bool`` to find the parts that are slug parameters and that need to be replaced. Data to change 
are taken from the real url by their numerical position.
