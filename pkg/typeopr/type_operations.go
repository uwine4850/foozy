package typeopr

import (
	"fmt"
	"reflect"
)

func IsPointer(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Pointer
}

func PtrIsStruct(a any) bool {
	return reflect.TypeOf(a).Elem().Kind() == reflect.Struct
}

type ErrValueNotPointer struct {
	Value string
}

func (e ErrValueNotPointer) Error() string {
	return fmt.Sprintf("%s value is not a pointer", e.Value)
}

type ErrValueIsPointer struct {
	Value string
}

func (e ErrValueIsPointer) Error() string {
	return fmt.Sprintf("%s value is a pointer", e.Value)
}

type ErrParameterNotStruct struct {
	Param string
}

func (e ErrParameterNotStruct) Error() string {
	return fmt.Sprintf("The %s parameter is not a structure.", e.Param)
}

type ErrConvertType struct {
	Type1 string
	Type2 string
}

func (e ErrConvertType) Error() string {
	return fmt.Sprintf("Data type conversion error. The %s interface type cannot be converted to %s type.", e.Type1, e.Type2)
}
