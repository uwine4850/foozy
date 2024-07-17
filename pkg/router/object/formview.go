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
	if err := form.FieldsNotEmpty(form.NewFillableFormStruct(&fillForm), v.NotNilFormFields); err != nil {
		return nil, err
	}

	return ObjectContext{namelib.OBJECT_CONTEXT_FORM: fillForm.Interface()}, nil
}
