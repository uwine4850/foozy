package object

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/form/formmapper"
	"github.com/uwine4850/foozy/pkg/secure"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type FormView struct {
	BaseView

	FormStruct       interface{} `notdef:"true"`
	NotNilFormFields []string
	NilIfNotExist    []string
	ValidateCSRF     bool
}

func (v *FormView) ObjectsName() []string {
	return []string{namelib.OBJECT.OBJECT_CONTEXT_FORM}
}

func (v *FormView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, managerConfig interfaces.IManagerConfig) (ObjectContext, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run FormView object", managerConfig)
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		return nil, err
	}
	if v.ValidateCSRF {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "validate CSRF token", managerConfig)
		if err := secure.ValidateFormCsrfToken(r, frm); err != nil {
			return nil, err
		}
	}

	debug.RequestLogginIfEnable(debug.P_OBJECT, "fill form", managerConfig)
	fillForm := reflect.New(reflect.TypeOf(v.FormStruct)).Elem()
	if err := formmapper.FillReflectValueFromForm(frm, &fillForm, v.NilIfNotExist); err != nil {
		return nil, err
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "check empty", managerConfig)
	if err := v.checkEmpty(typeopr.Ptr{}.New(&fillForm)); err != nil {
		return nil, err
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "check extension", managerConfig)
	if err := formmapper.CheckExtension(typeopr.Ptr{}.New(&fillForm)); err != nil {
		return nil, err
	}

	resultForm := fillForm.Interface()
	return ObjectContext{namelib.OBJECT.OBJECT_CONTEXT_FORM: &resultForm}, nil
}

// If the first character of the slice is "*", then you need to select the entire field of the structure.
// If there are more elements after the "*" sign, then they need to be excluded.
// When the "*" sign is missing, process according to the standard algorithm.
func (v *FormView) checkEmpty(fillForm itypeopr.IPtr) error {
	var notNilFields []string
	if len(v.NotNilFormFields) >= 1 && v.NotNilFormFields[0] == "*" {
		excludeFields := []string{}
		if len(v.NotNilFormFields) >= 2 {
			excludeFields = v.NotNilFormFields[1:]
		}
		fields, err := formmapper.FieldsName(fillForm, excludeFields)
		if err != nil {
			return err
		}
		notNilFields = fields
	} else {
		notNilFields = v.NotNilFormFields
	}
	if err := formmapper.FieldsNotEmpty(fillForm, notNilFields); err != nil {
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
