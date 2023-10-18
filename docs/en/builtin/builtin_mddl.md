## Package builtin_mddl
This package contains ready-made implementations of middleware.

__GenerateAndSetCsrf__
```
GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
This function is a standard middleware implementation.<br>
This function can be used to generate and set the value of csrf_token in the cookie parameters. Example of use:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
```