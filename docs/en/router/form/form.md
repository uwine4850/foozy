## Form package
The form package contains all available functionality (at the moment) for working with forms.

Form tests [here](https://github.com/uwine4850/foozy/tree/master/tests/form/form_test).

## Methods used by the IForm interface
__Parse__
```
Parse() error
```
The method that should always be used for form parsing. With its help, the basic processing of the form is performed.

__GetMultipartForm__
```
GetMultipartForm() *multipart.Form
```
Returns data arranged golang structure ``multipart.Form``. This method is used for the ``multipart/form-data`` form.

__GetApplicationForm__
```
GetApplicationForm() url.Values
```
Returns the data arranged golang structure ``url.Values``. This method is used for the ``application/x-www-form-urlencoded`` form.

__Value__
```
Value(key string) string
```
Returns data from the text field of the form.

__File__
```
File(key string) (multipart.File, *multipart.FileHeader, error)
```
Returns the form file data. Used with ``multipart/form-data`` form.

__Files__
```
Files(key string) ([]*multipart.FileHeader, bool)
```
Returns multiple files from the form (multiple input).

__ValidateCsrfToken__
```
ValidateCsrfToken() error
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

## Global package functions

__SaveFile__
```
SaveFile(w http.ResponseWriter, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string, manager interfaces.IManager) error
```
Saves the file in the selected location.
* fileHeader *multipart.FileHeader - information about the file.
* pathToDir string - the path to the directory to save the file.
* buildPath *string - reference to change string. The full path to the saved file is written in the change.

__ReplaceFile__
```
ReplaceFile(pathToFile string, w http.ResponseWriter, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string, manager interfaces.IManager) error
```
Replaces an existing file with another.

__SendApplicationForm__
```
SendApplicationForm(url string, values map[string]string) (*http.Response, error)
```
Sends a POST request (application/x-www-form-urlencoded) with data to the selected url.

__SendMultipartForm__
```
SendMultipartForm(url string, values map[string]string, files map[string][]string) (*http.Response, error)
```
Sends a POST request (multipart/form-data) with data to the selected url.
