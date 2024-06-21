## Package builtin_mddl
Цей пакет містить готові реалізації проміжного ПО.

__GenerateAndSetCsrf__
```
GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
This function is a standard implementation of middleware.<br>
With the help of this function, you can generate and set the csrf_token value in the cookie settings. Application example:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
```

__GenerateCsrfToken__
```
GenerateCsrfToken()
```
Generates a CSRF token.