## Form package
A package for working with HTML forms.<br>
The package performs the following tasks:

* Parsing forms `application/x-www-form-urlencoded` and `multipart/form-data`.
* Sending different types of forms.
* Working with files.

Example of use:
```golang
func Handler(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
    frm := form.NewForm(r)
    if err := frm.Parse(); err != nil {
        return err
    }
}
```

## FormFile
`FormFile` representation of the form file as an object.<br>
All data of the file is stored here.
```golang
type FormFile struct {
	Header *multipart.FileHeader
}
```

## Form
The `Form` object is designed to process forms. It parses `application/x-www-form-urlencoded` and `multipart/form-data` forms.

#### Form.Parse
Parses forms of types `application/x-www-form-urlencoded` and `multipart/form-data`, which are passed by `*http.Request`.
```golang
func (f *Form) Parse() error {
	enctype := strings.Split(f.request.Header.Get("Content-Type"), ";")[0]
	switch enctype {
	case "application/x-www-form-urlencoded":
		err := f.request.ParseForm()
		if err != nil {
			return err
		}
		f.applicationForm = f.request.PostForm
	case "multipart/form-data":
		err := f.request.ParseMultipartForm(f.multipartMaxMemory)
		if err != nil {
			return err
		}
		f.multipartForm = f.request.MultipartForm
	}
	return nil
}
```

#### Form.GetMultipartForm
Getting multipart/form-data form data.
```golang
func (f *Form) GetMultipartForm() *multipart.Form {
	return f.multipartForm
}
```

#### Form.GetApplicationForm
Getting multipart/form-data form data.
```golang
func (f *Form) GetApplicationForm() netUrl.Values {
	return f.applicationForm
}
```

#### Form.Value
Returns the value of the form by key. This can only be a simple value, a file cannot be returned here.
```golang
func (f *Form) Value(key string) string {
	val := f.request.PostFormValue(key)
	return template.HTMLEscapeString(val)
}
```

#### Form.File
Returns a file from a form by key.
```golang
func (f *Form) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return f.request.FormFile(key)
}
```

#### Form.File
Returns multiple files from the form by key.
```golang
func (f *Form) Files(key string) ([]*multipart.FileHeader, bool) {
	fi, ok := f.multipartForm.File[key]
	return fi, ok
}
```


#### SaveFile
SaveFile Saves the file in the specified directory.<br>
If the file name is already found, uses the randomiseTheFileName function to randomise the file name.
```golang
func SaveFile(fileHeader *multipart.FileHeader, pathToDir string, buildPath *string, manager interfaces.Manager) error
```

#### ReplaceFile
Changes the specified file to a new file.
```golang
func ReplaceFile(pathToFile string, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string, manager interfaces.Manager) error
```

#### SendApplicationForm
Sends a form of type application/x-www-form-urlencoded to the specified address.
```golang
func SendApplicationForm(url string, values map[string][]string) (*http.Response, error)
```

#### SendMultipartForm
Sends a form of type multipart/form-data to the specified address.<br>
The files argument accepts form field names and a slice with file paths.
```golang
func SendMultipartForm(url string, values map[string][]string, files map[string][]string) (*http.Response, error)
```