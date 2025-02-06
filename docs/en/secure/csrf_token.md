## package secure
This section of the security package is responsible for operations with the CSRF token.

__ValidateCookieCsrfToken__
```
ValidateCookieCsrfToken(r *http.Request, token string) error
```
The method that validates the CSRF token. To do this, you need to pass the token to the appropriate field.
The easiest way to do this is to add a built-in middleware [csrf](https://github.com/uwine4850/foozy/blob/master/docs/en/builtin/builtin_mddl/csrf.md) which will automatically add the ``csrf_token`` field to the cookie data. After which it can be validated here.
Also, you just need to add the variable ``{{ csrf_token | safe }}`` in the middle of the HTML form variable and run this method.<br>
Connecting the built-in middleware to create a token will happen as follows:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
...
newRouter.SetMiddleware(mddl)
```

__ValidateHeaderCSRFToken__
```
ValidateHeaderCSRFToken(r *http.Request, tokenName string)
```
The method that validates the CSRF token. For this it is necessary 
so that the token is transmitted through the header. You also need a token 
was in cookies before using this method.

__GenerateCsrfToken__
```
GenerateCsrfToken()
```
Generates a CSRF token.

__SetCSRFToken__
```
SetCSRFToken(maxAge int, w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error
```
Sets the token to a cookie. If __maxAge__ is zero, the cookie is session. 
Also adds the __CSRF_TOKEN__ variable to the templating context.