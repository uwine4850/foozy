## Fill struct

Form tests [here](https://github.com/uwine4850/foozy/tree/master/tests/formtest).

__type FillableFormStruct struct__

The FillableFormStruct structure is intended for more convenient access to the fillable structure.
The structure to be filled is passed as a pointer.

* _GetStruct() interface{}_ - returns the filled structure.<br>
* _SetDefaultValue(val func(name string) interface{})_ - sets the standard function.<br>
* _GetOrDef(name string, index int) interface{}_ - returns the structure value or a standard function if there is no structure value.<br>

__FillStructFromForm__
```
FillStructFromForm(frm *Form, fillableStruct *FillableFormStruct, nilIfNotExist []string) error
```
A method that fills a structure with data from a form.
A structure must always be passed as a reference.
For correct operation, it is necessary to specify the tag "form" for each field of the structure. For example, `form:<form field name>`.
* frm *Form - an instance of the form.
* fillableStruct *FillableFormStruct - instance of FillableFormStruct.
* nilIfNotExist - fields that are not found in the form will be nil.

__type OrderedForm struct__

A structure that organizes the form for further, more convenient use. All fields are ordered according to their order in the form.

* _Add(name string, value interface{})_ - adds a new form field to the structure.<br>
* _GetByName(name string) (OrderedFormValue, bool)_ - returns the form field by its name.<br>
* _GetAll() []OrderedFormValue_ - returns all form fields.<br>

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
Optimized to work even if the FillableFormStruct contains a structure with type *reflect.Value.

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
FormFile tag *ext* and expected extensions. For example, `ext:".jpeg .png"`.

__FillReflectValueFromForm__
```
FillReflectValueFromForm(frm *Form, fillValue *reflect.Value, nilIfNotExist []string) error
```
Fills the structure with data from the form.
The function works and does everything in the same way as the function `FillStructFromForm`.
The only difference is that this function accepts data in `*reflect.Value` format.