## Router package
The router module is an important module because it is the foundation of the entire web application. This module is divided into several parts,
namely:
* Router
* Manager
* Web sockets

## Router
The router is responsible for routing and working with handlers. Below are the methods that are available for use.<br>
### Route handlers
All handlers have standard parameters:
* *pattern* - the path to the handler. For example, the path can be "/home" and everyone is like it. The router also
    supports slug parameters, such as "/post/<id>". If you use such a path and go to the address "/post/1" - it will start
    the required handler in which the __id__ option will be available in the manager like this ``manager.GetSlugParams ("id")``. Such slug
    parameters can be many, the main thing with different names.
* *fn func(w http.ResponseWriter, r \*http.Request, manager interfaces.IManager)* - a function that will run when the user
  goes to the desired address. ``w http.ResponseWriter`` and ``*http.Request`` are standard golang structures. About ``interfaces.IManager``
  is described in more detail [here](#manager).
* Each handler returns ``func()`` - this is a function that is executed after the handler itself is finished.

__Get__
```
Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
The method is used to transfer data from the server to the web page. For example, it can be html data, JSON data, and others.

__Post__
```
Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
A method for processing a POST request. Most often used to work with HTML forms. This method has no new functionality
compared to the Get method. But this can be changed with the help of the [form handler package](https://github.com/uwine4850/foozy/blob/master/docs/en/form.md).

__Ws__
```
Ws(pattern string, ws interfaces2.IWebsocket, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager)) func()
```
The handler launches a web socket on the selected path. You can easily connect to this handler using JavaScript.
The ``interfaces2.IWebsocket`` parameter is the interface of a structure that implements interaction with a web socket, here is the standard implementation
``router.NewWebsocket(router.Upgrader)``. You can learn more about web sockets [here](#websocket).<br>
An example of implementing an echo handler:
```
newRouter.Ws("/ws", router.NewWebsocket(router.Upgrader), func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    ws := manager.GetWebSocket()
	err := ws.Connect(w, r, func() {
	    fmt.Println("Connect.")
	})
	if err != nil {
	    panic(err)
	}
	ws.OnClientClose(func() {
	    err := ws.Close()
		if err != nil {
		    panic(err)
		}
		fmt.Println("Client close.")
	})
	ws.OnMessage(func(messageType int, msgData []byte) {
	    err := ws.SendMessage(messageType, msgData)
		if err != nil {
		    panic(err)
		}
	})
	err = ws.ReceiveMessages()
	if err != nil {
	    panic(err)
	}
	return func(){}
})
```

__Put__
```
Put(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()
```
A method for handling PUT requests.

__Delete__
```
Delete(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()
```
Method for handling Delete requests.

### Other methods
__EnableLog__
```
EnableLog(enable bool)
```
Enable or disable logging to the console.

__SetMiddleware__
```
SetMiddleware(middleware interfaces.IMiddleware)
```
Sets an instance of the ``interfaces.IMiddleware`` interface to run before each handler.

__getHandleFunc__
```
getHandleFunc(pattern string, method string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager)) http.HandlerFunc
```
A private method that runs before each handler. It runs various validations, middleware, etc.
This method wraps the function ``func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)`` which is passed from
Post, Get or Ws handler.

__validateMethod__
```
validateMethod(method string) bool
```
Validation of http methods that are currently in use. That is, it makes sure that the Get method processes the http GET 
method, and not, for example, POST.

### Methods that do not belong to the interface but belong to the package
These methods are global, located in the router package, but can be used anywhere.<br>

__ValidateRootUrl__
```
ValidateRootUrl(w http.ResponseWriter, r *http.Request) bool
```
If the path pattern is __/__, then the handler will accept __all paths__. To prevent this, you need to use this
method to prevent this. Now if the path is not found, a 404 error will be displayed, not the __/__ handler. Example:
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
The function divides the url into parts by the character __/__. Then it adds each part to the map as a key in numerical order, and the value of the
of the key depends on whether it is a slug parameter, if it is, the key is true, if not, it is false.

__HandleSlugUrls__
```
HandleSlugUrls(parseUrl map[int]bool, slugUrl []string, url []string) (string, map[string]string)
```
``slugUrl []string`` is a pattern parameter that is separated by the __/__ character and written in slice.<br>
``url []string`` is a real url separated by __/__ and written in a slice.<br>

The function processes the url and outputs a string (str) as a url made from the pattern and the slug parameters if any.
Using the ``parseUrl map[int]bool``, you can find the parts that are slug parameters and need to be replaced. The data to be changed
are taken from the real url by their numeric position.

## Manager
This component is responsible for managing the handler. At this point, the interface performs the following functions:
* Transfer data from the middleware to the route handler.
* Replace the standard templating engine.
* Render the template.
* Receive slug parameters.
* Access to the web sockets interface.
* Render data in JSON format.

### Interface methods
__SetUserContext__
```
SetUserContext(key string, value interface{})
```
Sets the user context, for example, in middleware.

__GetUserContext__
```
GetUserContext(key string) (any, bool)
```
Returns the value of the user context by key.

__SetTemplateEngine__
```
SetTemplateEngine(engine interfaces2.ITemplateEngine)
```
Changes the standard templating tool.

__RenderTemplate__
```
RenderTemplate(w http.ResponseWriter, r *http.Request) error
```
Displays a template using a templating tool.<br>
__IMPORTANT:__ the template must be set using the ``SetTemplatePath`` method.

__SetTemplatePath__
```
SetTemplatePath(templatePath string)
```
Sets the path to the HTML template.

__SetContext__
```
SetContext(data map[string]interface{})
```
Sets the context for the templating agent. In an HTML template, this looks like ``{{ key }}``.

__SetSlugParams__
```
SetSlugParams(params map[string]string)
```
Sets the slug parameters. It is used in the router.

__GetSlugParams__
```
GetSlugParams(key string) (string, bool)
```
Provides access to slug parameters.

__SetWebsocket__
```
SetWebsocket(websocket interfaces.IWebsocket)
```
Sets up a web socket. It is used in the router.

__GetWebSocket__
```
GetWebSocket() interfaces.IWebsocket
```
Provides access to the web socket interface.

__RenderJson__
```
RenderJson(data interface{}, w http.ResponseWriter) error
```
Displays data in JSON format. As a parameter, data can take a map, structure, etc.

__DelUserContext__
```
DelUserContext(key string)
```
Deletes the user context by key.

## Websocket
The web socket interface is implemented using the __github.com/gorilla/websocket__ library. The ``router`` package has a global
variable ``Upgrader`` that is required for the web socket to work.

### Interface methods
__OnConnect__
```
OnConnect(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
```
Функція яка запускається під час підлючення до клієнта.

__Close__
```
Close() error
```
Closing the connection.

__OnClientClose__
```
OnClientClose(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
```
A function that will be executed when the client closes the connection.

__OnMessage__
```
OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
```
When the socket receives the message, the ``fn`` function will be executed.

__SendMessage__
```
SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
```
Sending a message to the client.

__ReceiveMessages__
```
ReceiveMessages(w http.ResponseWriter, r *http.Request) error
```
The method that starts receiving messages. This method must be running.