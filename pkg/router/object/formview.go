package object

import (
	"errors"
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
	ValidateCSRF     bool
}

func (v *FormView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
	if typeopr.IsPointer(v.FormStruct) {
		return nil, typeopr.ErrValueIsPointer{Value: "FormStruct"}
	}
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		return nil, err
	}
	if v.ValidateCSRF {
		if err := frm.ValidateCsrfToken(); err != nil {
			return nil, err
		}
	}

	fillForm := reflect.New(reflect.TypeOf(v.FormStruct)).Elem()
	if err := form.FillReflectValueFromForm(frm, &fillForm, v.NilIfNotExist); err != nil {
		return nil, err
	}
	fillableForm := form.NewFillableFormStruct(&fillForm)

	if err := v.checkEmpty(fillableForm); err != nil {
		return nil, err
	}

	if err := form.CheckExtension(fillableForm); err != nil {
		return nil, err
	}

	resultForm := fillForm.Interface()
	return ObjectContext{namelib.OBJECT.OBJECT_CONTEXT_FORM: &resultForm}, nil
}

// If the first character of the slice is "*", then you need to select the entire field of the structure.
// If there are more elements after the "*" sign, then they need to be excluded.
// When the "*" sign is missing, process according to the standard algorithm.
func (v *FormView) checkEmpty(fillableForm *form.FillableFormStruct) error {
	var notNilFields []string
	if len(v.NotNilFormFields) >= 1 && v.NotNilFormFields[0] == "*" {
		excludeFields := []string{}
		if len(v.NotNilFormFields) >= 2 {
			excludeFields = v.NotNilFormFields[1:]
		}
		fields, err := form.FieldsName(fillableForm, excludeFields)
		if err != nil {
			return err
		}
		notNilFields = fields
	} else {
		notNilFields = v.NotNilFormFields
	}
	if err := form.FieldsNotEmpty(fillableForm, notNilFields); err != nil {
		return err
	}
	return nil
}

// FormInterface retrieves the form interface itself from the interface pointer.
func (v *FormView) FormInterface(manager interfaces.IManagerOneTimeData) (interface{}, error) {
	context, ok := manager.GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	if !ok {
		return nil, errors.New("the ObjectContext not found")
	}
	objectContext, ok := context.(ObjectContext)
	if !ok {
		return nil, errors.New("the ObjectContext type assertion error")
	}
	return reflect.Indirect(reflect.ValueOf(objectContext[namelib.OBJECT.OBJECT_CONTEXT_FORM])).Interface(), nil
}
