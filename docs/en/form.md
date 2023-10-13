## Form package
The form package contains all the available functionality (at the moment) for working with forms.

## Methods used by the IForm interface
__Parse__
```
Parse() error
```
A method that should always be used for form parsing. It is used to perform basic form processing.

__GetMultipartForm__
```
GetMultipartForm() *multipart.Form
```
Returns the data of the organized golang structure ``multipart.Form``. This method is used for the ``multipart/form-data`` form.

__GetApplicationForm__
```
GetApplicationForm() url.Values
```
Returns a golang ``url.Values`` structure. This method is used for the ``application/x-www-form-urlencoded`` form.

__Value__
```
Value(key string) string
```
Returns data from the form text field.

__File__
```
File(key string) (multipart.File, *multipart.FileHeader, error)
```
Returns the data of the form file. Used with the ``multipart/form-data`` form.

__ValidateCsrfToken__
```
ValidateCsrfToken() error
```
A method that validates the CSRF token. To do this, the form must have a field called ``csrf_token``, in addition, the cookie data
must also have a field called ``csrf_token``.<br>
The easiest way to do this is to add built-in [middleware](https://github.com/uwine4850/foozy/blob/master/docs/en/middlewares.md)
that will automatically add the ``csrf_token`` field to the cookie data.
After that, you just need to add the variable ``{{ csrf_token }}`` to the middle of the HTML form and run this method.<br>
Connect the built-in middleware to create a token as follows:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
...
newRouter.SetMiddleware(mddl)
```
