package formmapper

import (
	"mime/multipart"
	"net/url"

	"github.com/uwine4850/foozy/pkg/router/form"
)

type IFormGetEnctypeData interface {
	GetMultipartForm() *multipart.Form
	GetApplicationForm() url.Values
}

// OrderedForm Values can be displayed either by field name or all fields at once.
type OrderedForm struct {
	itemCount int
	names     map[string][]int
	values    []OrderedFormValue
}

func NewOrderedForm() *OrderedForm {
	o := &OrderedForm{}
	o.itemCount = 0
	o.names = make(map[string][]int)
	return o
}

// Add a new form field.
func (f *OrderedForm) Add(name string, value interface{}) {
	f.values = append(f.values, OrderedFormValue{
		Name:  name,
		Value: value,
	})
	f.itemCount++
	f.names[name] = append(f.names[name], f.itemCount)
}

// GetByName getting a field by name.
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

// GetAll getting all fields.
func (f *OrderedForm) GetAll() []OrderedFormValue {
	return f.values
}

type OrderedFormValue struct {
	Name  string
	Value interface{}
}

// FrmValueToOrderedForm Converts the form to a OrderedForm.
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
