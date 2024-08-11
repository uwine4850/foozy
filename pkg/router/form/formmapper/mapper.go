package formmapper

import (
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/router/form"
)

type Mapper struct {
	Form          *form.Form
	Output        itypeopr.IPtr
	NilIfNotExist []string
}

func NewMapper(form *form.Form, output itypeopr.IPtr, nilIfNotExist []string) Mapper {
	return Mapper{Form: form, Output: output, NilIfNotExist: nilIfNotExist}
}

func (m *Mapper) Fill() error {
	switch reflect.TypeOf(m.Output.Ptr()).Elem().Kind() {
	case reflect.Struct:
		// *reflect.Value
		if reflect.DeepEqual(reflect.TypeOf(&reflect.Value{}), reflect.TypeOf(m.Output.Ptr())) {
			if err := FillReflectValueFromForm(m.Form, m.Output.Ptr().(*reflect.Value), m.NilIfNotExist); err != nil {
				return err
			}
			if err := CheckExtension(m.Output); err != nil {
				return err
			}
		} else {
			// struct
			if err := FillStructFromForm(m.Form, m.Output, m.NilIfNotExist); err != nil {
				return err
			}
			if err := CheckExtension(m.Output); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("mapping for %s type is not supported", reflect.TypeOf(m.Output.Ptr()).Kind())
	}
	return nil
}
