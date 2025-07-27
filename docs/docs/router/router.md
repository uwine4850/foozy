## package router

Globally, the router package is responsible for two functions:

1. Preparing and running handlers via [Adapter](#adapter).
2. Registering and running http methods. This is the responsibility of [Router](#router).

In general, the algorithm of route processing is as follows:

1. Register the handler on the route using the [Register](#routerregister) method.
2. Passing the handler to [Adapter](#adapter) and receiving the same handler, but in a special wrapper.
3. Adding the http method, path pattern and adapted handler to the route processing list.
4. When url and path pattern match, the adapted handler will start.

### Adapter
The `Adapter` object is needed to adapt the handler to work with the rest of the framework modules. The object is also needed to control the handler in a more flexible way. In general, `Adapter` should wrap the passed handler to work with other modules, and then this wrapper will be launched in [Router](#router).

The adapter's responsibilities include the following items:

* Clearing logs before running the handler
* Creating a new instance of the [Manager](/router/manager/manager/) object.
* Setting some [Manager.OneTimeData](/router/manager/manager/#onetimedata) variables
* Running all types of middleware
* Running the original handler
* Outputting request logs to the console

#### Adapter.Adapt
This is the main method that wraps the source handler and thus runs all the necessary modules. In general, its implementation is quite simple, but there are a few nuances:

* Instead of `http.ResponseWriter`, [BufferedResponseWriter](#bufferedresponsewriter) is used.
* If the connection is a websocket, [PostMiddlewares](/router/middlewares/middlewares/#middlewarespostmiddleware) are not run

### BufferedResponseWriter
BufferedResponseWriter is a wrapper over `http.ResponseWriter`. The wrapper is needed to give more flexible control over how data is written to the page. This object fully implements the `http.ResponseWriter` interface, but the difference is that the `Write` method does not send the request immediately, but buffers it. To actually write the data you need to use the [Flush](#bufferedresponsewriterflush) method.

#### BufferedResponseWriter.Flush
A simple method that writes data to the original `http.ResponseWriter`. Accordingly, the data is sent immediately.
```golang
func (rw *BufferedResponseWriter) Flush() (int, error) {
	for k, vv := range rw.header {
		for _, v := range vv {
			rw.original.Header().Add(k, v)
		}
	}
	rw.original.WriteHeader(rw.statusCode)
	return rw.original.Write(rw.buffer.Bytes())
}
```

### Router
The `Router` object is used to route http requests. The algorithm of its work looks like this:

1. Obtaining the url address
2. Searching for url address match in stored url templates
3. Running a handler that is bound to the url template

#### Router.HandlerSet
Gegisters multiple handlers at once. Does everything the same as [Router.Register](#routerregister), but only with multiple handlers. The method is just for convenience.

Example Usage:
```golang
newRouter.HandlerSet(
    map[string][]map[string]router.Handler{
		router.MethodGET: {
			{"/home": homeHandler},
			{"/about": aboutHandler},
		},
		router.MethodPOST: {
			{"/register": registerHandler},
		},
	}
)
```

#### Router.Register
It simply registers a handler on the selected path. Uses [Adapter](router.md#adapter) to wrap the handler.
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

#### Router.ServeHTTP
Implements the `http.Handler` interface. It is used to call handlers.


#### RedirectError
The function redirects to the selected url. Also, the url parameters are redirected with an error from the function arguments.
```golang
func RedirectError(w http.ResponseWriter, r *http.Request, path string, _err string) {
	uval := url.Values{}
	uval.Add(namelib.ROUTER.REDIRECT_ERROR, _err)
	newUrl := fmt.Sprintf("%s?%s", path, uval.Encode())
	http.Redirect(w, r, newUrl, http.StatusFound)
	debug.ErrorLogginIfEnable(_err)
	debug.RequestLogginIfEnable(debug.P_ERROR, _err)
}
```

#### CatchRedirectError
Catches an error from url parameters that is passed through the [RedirectError](#redirecterror) function. Sets the error to the [manager](/router/manager/manager/) context. If there is a __TODO: link__ [Render]() instance in [Manager](/router/manager/manager/), sets the error to its context. Both contexts are handled by the `namelib.ROUTER.REDIRECT_ERROR` key.
```golang
func CatchRedirectError(r *http.Request, manager interfaces.Manager) {
	q := r.URL.Query()
	redirectError := q.Get(namelib.ROUTER.REDIRECT_ERROR)
	if redirectError != "" {
		if manager.Render() != nil {
			manager.Render().SetContext(map[string]interface{}{namelib.ROUTER.REDIRECT_ERROR: redirectError})
		}
		manager.OneTimeData().SetUserContext(namelib.ROUTER.REDIRECT_ERROR, redirectError)
	}
}
```

#### ServerError
Sends an error with code 500 and the text "500 Internal server error" to the page. The text is sent using the special function __TODO: link__ [debug.ErrorLoggingIfEnableAndWrite]().
```golang
func ServerError(w http.ResponseWriter, error string, manager interfaces.Manager) {
	manager.OneTimeData().SetUserContext(namelib.ROUTER.SERVER_ERROR, error)
	w.WriteHeader(http.StatusInternalServerError)
	if config.LoadedConfig().Default.Debug.Debug {
		debug.ErrorLoggingIfEnableAndWrite(w, error, error)
	} else {
		debug.ErrorLoggingIfEnableAndWrite(w, error, "500 Internal server error")
	}
	debug.RequestLogginIfEnable(debug.P_ERROR, error)
}
```

#### ServerForbidden
Sends an error with code 403 and text "500 Internal server error" to the page. The text is sent using the special function __TODO: link__ [debug.ErrorLoggingIfEnableAndWrite]().
```golang
func ServerForbidden(w http.ResponseWriter, manager interfaces.Manager) {
	manager.OneTimeData().SetUserContext(namelib.ROUTER.SERVER_FORBIDDEN_ERROR, "403 forbidden")
	w.WriteHeader(http.StatusForbidden)
	debug.ErrorLoggingIfEnableAndWrite(w, "403 forbidden", "403 forbidden")
	debug.RequestLogginIfEnable(debug.P_ERROR, "403 forbidden")
}
```

#### SendJson
Sends Json as a response to an http request.
```golang
func SendJson(data interface{}, w http.ResponseWriter, code int) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}
```