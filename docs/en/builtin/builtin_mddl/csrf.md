## Package builtin_mddl

__GenerateAndSetCsrf__
```
GenerateAndSetCsrf(maxAge int, onError onError) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
This function returns a standard middleware implementation.<br>
maxAge - cookie lifetime.
onError - a function that will be executed during an error.
With the help of this function, you can generate and set the csrf_token value in the cookie settings. Application example:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
```