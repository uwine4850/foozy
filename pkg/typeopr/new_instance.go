package typeopr

import (
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces/intrnew"
)

// CreateNewInstance сreates a new instance of a structure. The structure must implement the interface interfaces.INewInstance.
// The <new> argument takes a pointer to the structure that will contain the new instance.
func CreateNewInstance(ins intrnew.INewInstance, new interface{}) error {
	if !IsPointer(new) {
		panic(ErrValueNotPointer{Value: "new"})
	}
	var typeIns reflect.Type
	if reflect.TypeOf(ins).Kind() == reflect.Ptr {
		typeIns = reflect.TypeOf(ins).Elem()
	} else {
		typeIns = reflect.TypeOf(ins)
	}

	reflectIns := reflect.New(typeIns).Interface().(intrnew.INewInstance)
	newIns, err := reflectIns.New()
	if err != nil {
		return err
	}
	newInsInterface := reflect.ValueOf(newIns).Interface()
	reflect.ValueOf(new).Elem().Set(reflect.ValueOf(newInsInterface))
	return nil
}
