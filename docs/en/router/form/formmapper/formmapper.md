## package formmapper
Fills the structure with data from the form.

You can see how the package works in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/formtest/formmapping_test).

### type Mapper struct

`Form` — a reference to an already parsed form.<br>
`Output` — a reference to a structure, or a structure in `*reflect.Value`.<br>
`NilIfNotExist` — if no matching fields are found in the form and they are in this list, the structure fields will be nil.

An example of configuring the structure:
```
type Fill struct {
	Field1   []string        `form:"f1"   empty:"1"`
	Field2   []string        `form:"f2"   empty:"-err"`
	File     []form.FormFile `form:"file" empty:"-err"`
}
```
Each structure field must be of type `[]string` or `[]form.FormFile`. This is necessary in order to 
so that one name (key) input can be used to transfer several calculations.<br>

The `form` tag is required. It is responsible for the name of the input, which will be written in the field of the structure.<br>

The `empty` tag will only be applied to empty values. It is important to note here that this only applies 
for the slice value, that is, if there are two values ​​for the key, but it is empty 
only one, then the operation will be performed only with an empty value. This tag has 
several options:
*  -err — throws an error if at least one of the slice indices is empty. 
The type `[]form.FormFile` can have only one option - `-err`.
*  plain text — will replace the data of the empty value by its index.

The `ext` tag is an extension of the files that can be in the field. For example `ext:".jpg .png"`.

__Fill__
```
Fill() error
```
Fills the structure with values ​​from the form.

### type OrderedForm struct

A structure that organizes the form for further, more convenient use. All fields are ordered according to their order in the form.

* _Add(name string, value interface{})_ - adds a new form field to the structure.<br>
* _GetByName(name string) (OrderedFormValue, bool)_ - returns a form field by its name.<br>
* _GetAll() []OrderedFormValue_ - returns all form fields.<br>
  
### Other features of the package

__FillStructFromForm__
```
FillStructFromForm(frm *Form, fillStruct interface{}, nilIfNotExist []string) error
```
A method that fills a structure with data from a form.
A structure must always be passed as a reference.
For correct operation, it is necessary to specify the tag "form" for each field of the structure. For example, `form:<form field name>`. Also supports the `empty` tag described above.
* frm *Form - form instance.
* fillStruct interface{} - reference to the object to be filled.
* nilIfNotExist - fields that are not found in the form will be nil.

__FrmValueToOrderedForm__
```
FrmValueToOrderedForm(frm IFormGetEnctypeData) *OrderedForm
```
Fills the form data into the *OrderedForm* structure.

__FieldsNotEmpty__
```
FieldsNotEmpty(fillStruct interface{}, fieldsName []string) error
```
Checks whether the selected structure fields are empty.
Optimized to work even if FillableFormStruct contains a structure with type *reflect.Value.

__FieldsName__
```
FieldsName(fillStruct interface{}, exclude []string) ([]string, error)
```
Returns the names of the structure's fields.

__CheckExtension__
```
CheckExtension(fillStruct interface{}) error
```
Checks if the extension of the form files is as expected. To work correctly, you need to add a type to each field 
FormFile tag *ext* and expected extensions. For example, `ext:".jpeg .png"`.

__FillReflectValueFromForm__
```
FillReflectValueFromForm(frm *Form, fillValue *reflect.Value, nilIfNotExist []string) error
```
Fills the structure with data from the form.
The function works and does everything in the same way as the function `FillStructFromForm`.
The only difference is that this function accepts data in `*reflect.Value` format.