package formmapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type Mapper struct {
	Form          *form.Form
	Output        interface{}
	NilIfNotExist []string
}

func NewMapper(form *form.Form, output interface{}, nilIfNotExist []string) Mapper {
	return Mapper{Form: form, Output: output, NilIfNotExist: nilIfNotExist}
}

func (m *Mapper) Fill() error {
	outType, err := m.outputType()
	if err != nil {
		return err
	}
	switch outType {
	case reflect.Struct:
		// *reflect.Value
		if reflect.DeepEqual(reflect.TypeOf(&reflect.Value{}), reflect.TypeOf(m.Output)) {
			if err := FillReflectValueFromForm(m.Form, m.Output.(*reflect.Value), m.NilIfNotExist); err != nil {
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
		return fmt.Errorf("mapping for %s type is not supported", outType)
	}
	return nil
}

func (m *Mapper) outputType() (reflect.Kind, error) {
	if !typeopr.IsPointer(m.Output) {
		return reflect.Invalid, typeopr.ErrValueNotPointer{Value: "Output"}
	}
	typeOf := reflect.TypeOf(m.Output).Elem()
	if typeOf.Kind() != reflect.Struct {
		return reflect.Invalid, errors.New("field Output must be a struct")
	} else {
		return typeOf.Kind(), nil
	}
}
