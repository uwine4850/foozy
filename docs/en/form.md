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
After that, you just need to add the variable ``{{ csrf_token | safe }}`` to the middle of the HTML form and run this method.<br>
Connect the built-in middleware to create a token as follows:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
...
newRouter.SetMiddleware(mddl)
```

__Files__
```
Files(key string) ([]*multipart.FileHeader, bool)
```
Returns multiple files from the form(multiple input).

__SaveFile__
```
SaveFile(w http.ResponseWriter, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) error
```
Saves the file in the selected location.
* fileHeader *multipart.FileHeader - file information.
* pathToDir string - path to the directory to save the file.
* buildPath *string - request to change string. The full path to the saved file is written in the change.

__ReplaceFile__
```
ReplaceFile(pathToFile string, w http.ResponseWriter, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) error
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

### Fill struct

__type FillableFormStruct struct__

The FillableFormStruct structure is intended for more convenient access to the fillable structure.
The structure to be filled is passed as a pointer.

* _GetStruct() interface{}_ - Returns the filled structure.<br>
* _SetDefaultValue(val func(name string) string)_ - sets the standard function.<br>
* _GetOrDef(name string, index int) string_ - returns the structure value or a standard function if the structure value is missing.<br>


## Global functions of the package
__FillStructFromForm__.
```
FillStructFromForm(frm *Form, fillableStruct *FillableFormStruct, nilIfNotExist []string) error
```
A method that fills a structure with data from a form.
The structure must always be passed as a reference.
For correct operation, you must specify the "form" tag for each field of the structure. For example, `form:<form field name>`.
* frm *Form - form instance.
* fillableStruct *FillableFormStruct - an instance of FillableFormStruct.
* nilIfNotExist - fields that are not found in the form will be nil.

__type OrderedForm struct__

A structure that organizes the form for later more convenient use. All fields are ordered according to their order in the form.

* _Add(name string, value interface{})_ - adds a new form field to the structure.<br>
* GetByName(name string) (OrderedFormValue, bool)_ - returns a form field by its name.<br>
* GetAll() []OrderedFormValue_ - returns all form fields.<br>

__FrmValueToOrderedForm__
```
FrmValueToOrderedForm(frm IFormGetEnctypeData) *OrderedForm
```
Fills the form data into the *OrderedForm* structure.

__FieldsNotEmpty__
```
FieldsNotEmpty(fillableStruct *FillableFormStruct, fieldsName []string) error
```
Checks whether the selected structure fields are empty.

__FieldsName__
```
FieldsName(fillForm *FillableFormStruct, exclude []string) ([]string, error)
```
Returns the names of the structure's fields.

__CheckExtension__
```
CheckExtension(fillForm *FillableFormStruct) error
```
Checks if the extension of the form files is as expected. To work correctly, you need to add a type to each field
FormFile tag *ext* and expected extensions. For example, `ext:".jpeg, .png"`.
