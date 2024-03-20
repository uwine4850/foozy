## Middlewares package
This package contains all the tools you need to work with middleware.
## Методи які використовує інтерфейс IMiddleware
Methods that create a handler have one common parameter ``fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData)`` - 
this parameter is identical to the parameter described in the [router package](https://github.com/uwine4850/foozy/blob/master/docs/ua/router.md).<br>
The only difference is that the data of this parameter first goes to middlewares, and then to the router handler.

__HandlerMddl__
```
HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData))
```
The method creates a middleware that will be executed synchronously. The ``id`` parameter is the sequence number of the middleware execution, it
must be unique.

__AsyncHandlerMddl__
```
AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData))
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

### Методи пакету
__SetMddlError__
```
SetMddlError(mddlErr error, manager interfaces.IManagerData)
```
Saves the error that occurred in the middleware.

__GetMddlError__
```
GetMddlError(manager interfaces.IManagerData) (error, error)
```
Returns the error that was set using __SetMddlError__.

__SkipNextPage__
```
SkipNextPage(manager interfaces.IManagerData)
```
Skips the rendering of the next web request. The handler does not start even at the initial stage.

__IsSkipNextPage__
```
IsSkipNextPage(manager interfaces.IManagerData) bool
```
Checks whether the next page should be skipped. In the standard implementation, it is used in a router.

__SkipNextPageAndRedirect__
```
SkipNextPageAndRedirect(manager interfaces.IManagerData, w http.ResponseWriter, r *http.Request, path string)
```
Skips the rendering of the next page and redirects to another.
