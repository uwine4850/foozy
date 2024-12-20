package fstruct

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

// CheckNotDefaultFields checks whether the values ​​of the structure fields are default.
// That is, if the field is not passed or initialized with standard values, for example, nil.
func CheckNotDefaultFields(objectPtr typeopr.IPtr) error {
	objectLink := objectPtr.Ptr()
	var rObjectValue reflect.Value
	var rObjectType reflect.Type
	if reflect.TypeOf(objectLink).Elem() == reflect.TypeOf(reflect.Value{}) {
		rObjectValue = objectLink.(*reflect.Value).Elem()
		rObjectType = objectLink.(*reflect.Value).Elem().Type()
	} else {
		rObjectValue = reflect.ValueOf(objectLink).Elem()
		rObjectType = reflect.TypeOf(objectLink).Elem()
	}
	for i := 0; i < rObjectType.NumField(); i++ {
		fieldType := rObjectType.Field(i)
		fieldValue := rObjectValue.Field(i)
		tag := fieldType.Tag.Get("notdef")
		if tag == "" {
			continue
		}
		reqiredValue, err := strconv.ParseBool(tag)
		if err != nil {
			return err
		}
		if reqiredValue {
			if fieldValue.IsZero() {
				return ErrStructFieldIsDefault{fieldType.Name}
			}
		}
	}
	return nil
}

type ErrStructFieldIsDefault struct {
	FieldName string
}

func (e ErrStructFieldIsDefault) Error() string {
	return fmt.Sprintf("struct field %s is default", e.FieldName)
}
