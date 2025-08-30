## orderedform

### OrderedForm
This structure is used to store the ordered data of the form.<br>
This is necessary to ensure that the data is in the same order each time the form is used.
```golang
type OrderedForm struct {
	itemCount int
	names     map[string][]int
	values    []OrderedFormValue
}
```

#### OrderedForm.Add
Add a new form field.
```golang
func (f *OrderedForm) Add(name string, value interface{}) {
	f.values = append(f.values, OrderedFormValue{
		Name:  name,
		Value: value,
	})
	f.itemCount++
	f.names[name] = append(f.names[name], f.itemCount)
}
```

#### OrderedForm.GetByName
Getting a field by name.
```golang
func (f *OrderedForm) GetByName(name string) ([]OrderedFormValue, bool) {
	getIndex, ok := f.names[name]
	if !ok {
		return nil, ok
	}
	res := []OrderedFormValue{}
	for i := 0; i < len(getIndex); i++ {
		res = append(res, f.values[getIndex[i]-1])
	}
	return res, true
}
```

#### OrderedForm.GetAll
Getting all fields.
```golang
func (f *OrderedForm) GetAll() []OrderedFormValue {
	return f.values
}
```

### OrderedFormValue
Single value of a form field.
```golang
type OrderedFormValue struct {
	Name  string
	Value interface{}
}
```

#### FrmValueToOrderedForm
Converts the form to a `OrderedForm`.
```golang
func FrmValueToOrderedForm(frm IFormGetEnctypeData) *OrderedForm {
	orderedForm := NewOrderedForm()
	multipartForm := frm.GetMultipartForm()
	if multipartForm != nil {
		for name, value := range multipartForm.Value {
			orderedForm.Add(name, value)
		}
		for name, value := range multipartForm.File {
			var files []form.FormFile
			for i := 0; i < len(value); i++ {
				files = append(files, form.FormFile{Header: value[i]})
			}
			orderedForm.Add(name, files)
		}
	}
	applicationForm := frm.GetApplicationForm()
	if applicationForm != nil {
		for name, value := range applicationForm {
			orderedForm.Add(name, value)
		}
	}
	return orderedForm
}
```