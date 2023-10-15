## Middlewares package
This package contains all the tools you need to work with middleware.
## Методи які використовує інтерфейс IMiddleware
Methods that create a handler have one common parameter ``fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)`` - 
this parameter is identical to the parameter described in the [router package](https://github.com/uwine4850/foozy/blob/master/docs/ua/router.md).<br>
The only difference is that the data of this parameter first goes to middlewares, and then to the router handler.

__HandlerMddl__
```
HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager))
```
The method creates a middleware that will be executed synchronously. The ``id`` parameter is the sequence number of the middleware execution, it
must be unique.

__AsyncHandlerMddl__
```
AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager))
```
This method can also be used to create middleware, but it will run asynchronously. Accordingly, there are no sequence numbers.
__RunMddl__
```
RunMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error
```
Launching synchronous middleware.

__RunAsyncMddl__
```
RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
Launching asynchronous middleware.<br>
__IMPORTANT__: all middleware runs asynchronously, so you need to wait for their execution using the ``WaitAsyncMddl`` method.

__WaitAsyncMddl__
```
WaitAsyncMddl()
```
Wait for the execution of all asynchronous middleware (if any).
