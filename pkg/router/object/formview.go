package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type FormView struct {
	BaseView

	FormStruct       interface{}
	NotNilFormFields []string
	NilIfNotExist    []string
}

func (v *FormView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
	if typeopr.IsPointer(v.FormStruct) {
		panic("the FormStruct parameter must not be a pointer")
	}
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		return nil, err
	}

	fillForm := reflect.New(reflect.TypeOf(v.FormStruct)).Elem()
	if err := form.FillReflectValueFromForm(frm, &fillForm, v.NilIfNotExist); err != nil {
		return nil, err
	}

	// If the first character of the slice is "*", then you need to select the entire field of the structure.
	// If there are more elements after the "*" sign, then they need to be excluded.
	// When the "*" sign is missing, process according to the standard algorithm.
	fillableForm := form.NewFillableFormStruct(&fillForm)
	var notNilFields []string
	if len(v.NotNilFormFields) >= 1 && v.NotNilFormFields[0] == "*" {
		excludeFields := []string{}
		if len(v.NotNilFormFields) >= 2 {
			excludeFields = v.NotNilFormFields[1:]
		}
		fields, err := form.FieldsName(fillableForm, excludeFields)
		if err != nil {
			return nil, err
		}
		notNilFields = fields
	} else {
		notNilFields = v.NotNilFormFields
	}
	if err := form.FieldsNotEmpty(fillableForm, notNilFields); err != nil {
		return nil, err
	}

	return ObjectContext{namelib.OBJECT_CONTEXT_FORM: fillForm.Interface()}, nil
}
