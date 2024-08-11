package typeopr

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

func IsPointer(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Pointer
}

func PtrIsStruct(a any) bool {
	return reflect.TypeOf(a).Elem().Kind() == reflect.Struct
}

func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		return v.NumField() == 0
	case reflect.Invalid:
		return true
	}
	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}

func AnyToBytes(value interface{}) ([]byte, error) {
	var buf bytes.Buffer

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("%v", val.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(fmt.Sprintf("%v", val.Uint()))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("%v", val.Float()))
	case reflect.String:
		buf.WriteString(val.String())
	default:
		return nil, fmt.Errorf("unsupported convert type %s", val.Kind().String())
	}
	return buf.Bytes(), nil
}

type Ptr struct {
	value interface{}
}

func (p Ptr) New(value interface{}) itypeopr.IPtr {
	if !IsPointer(value) {
		panic(ErrValueNotPointer{Value: fmt.Sprintf("Ptr<%s>", reflect.TypeOf(value))})
	}
	p.value = value
	return p
}

func (p Ptr) Ptr() interface{} {
	return p.value
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

type ErrValueIsEmpty struct {
	Value string
}

func (e ErrValueIsEmpty) Error() string {
	return fmt.Sprintf("value %s is empty", e.Value)
}
