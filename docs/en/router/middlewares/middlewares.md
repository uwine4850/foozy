## Middlewares package
This package contains all the necessary tools for working with middleware.

## Methods used by the IMiddleware interface
Methods that create a handler have one common parameter ``fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData)`` - 
this parameter is identical to the parameter described in the package [router](https://github.com/uwine4850/foozy/blob/master/docs/en/router/router.md).<br>
The only difference is that the data of this parameter first goes to the middlewares, and then to the router processor.

Tests for Middlewares [here](https://github.com/uwine4850/foozy/tree/master/tests/middlewares).

__HandlerMddl__
```
HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData))
```
The method creates middleware that will be executed synchronously. The ``id`` parameter is the serial number of the middleware execution, it 
must be unique.

__AsyncHandlerMddl__
```
AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData))
```
You can also create a middleware using this method, but it will run asynchronously. Accordingly, serial numbers do not exist.

__RunMddl__
```
RunMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error
```
Running synchronous middleware.

__RunAsyncMddl__
```
RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
Running asynchronous middleware.<br>
__IMPORTANT__: All middlewares run asynchronously, so you need to wait for them to complete using the ``WaitAsyncMddl`` method.

__WaitAsyncMddl__
```
WaitAsyncMddl()
```
Wait for the execution of all asynchronous middleware (if any).

### Package methods
__SetMddlError__
```
SetMddlError(mddlErr error, manager interfaces.IManagerOneTimeData)
```
Saves the error that occurred in the middleware.

__GetMddlError__
```
GetMddlError(manager interfaces.IManagerOneTimeData) (error, error)
```
Returns the error that was set using __SetMddlError__.

__SkipNextPage__
```
SkipNextPage(manager interfaces.IManagerOneTimeData)
```
Skips the rendering of the next web request. The handler does not start even at the initial stage.

__IsSkipNextPage__
```
IsSkipNextPage(manager interfaces.IManagerOneTimeData) bool
```
Checks whether the next page should be skipped. In the standard implementation, it is used in a router.

__SkipNextPageAndRedirect__
```
SkipNextPageAndRedirect(manager interfaces.IManagerOneTimeData, w http.ResponseWriter, r *http.Request, path string)
```
Skips the rendering of the next page and redirects to another.
