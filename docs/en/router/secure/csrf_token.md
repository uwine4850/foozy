## package secure
This section of the security package is responsible for operations with the CSRF token.

__ValidateFormCsrfToken__
```
ValidateFormCsrfToken(r *http.Request, frm *form.Form) error
```
The method that validates the CSRF token. For this, the form must have a field called ``csrf_token``, in addition to cookies
must also have a ``csrf_token`` field.<br>
The easiest way to do this is to add a built-in middleware [csrf](https://github.com/uwine4850/foozy/blob/master/docs/en/builtin/builtin_mddl/csrf.md) which will automatically add the ``csrf_token`` field to the cookie data.
After that, you just need to add the variable ``{{ csrf_token | safe }}`` in the middle of the HTML form variable and run this method.<br>
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