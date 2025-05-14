package typeopr

import (
	"bytes"
	"fmt"
	"reflect"
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

type IPtr interface {
	New(value interface{}) IPtr
	Ptr() interface{}
}

type Ptr struct {
	value interface{}
}

func (p Ptr) New(value interface{}) IPtr {
	if !IsPointer(value) {
		panic(ErrValueNotPointer{Value: fmt.Sprintf("Ptr<%s>", reflect.TypeOf(value))})
	}
	p.value = value
	return p
}

func (p Ptr) Ptr() interface{} {
	return p.value
}

// IsImplementInterface determines whether an object uses an interface.
// Usage example:
// object := MyObject{}
// IsImplementInterface(typeopr.Ptr{}.New(&object), (*MyInterface)(nil))
// If reflect.Value is used, you can use direct passing or passing by pointer, that is,
// passing a pointer to a pointer. How to transmit this data depends on the situation.
func IsImplementInterface(objectPtr IPtr, interfaceType interface{}) bool {
	object := objectPtr.Ptr()
	// If the type of data passed directly is the desired interface.
	if reflect.TypeOf(object) == reflect.TypeOf(interfaceType) {
		return true
	}
	var objType reflect.Type
	if reflect.TypeOf(object).Elem() == reflect.TypeOf(reflect.Value{}) {
		objType = object.(*reflect.Value).Type()
	} else {
		objType = reflect.TypeOf(object)
	}
	intrfcType := reflect.TypeOf(interfaceType).Elem()
	return objType.Implements(intrfcType)
}

var typeReflectValue = reflect.TypeOf(&reflect.Value{}).Elem()

func GetReflectValue[T any](target *T) reflect.Value {
	var v reflect.Value
	if reflect.TypeOf(*target) == typeReflectValue {
		v = *reflect.ValueOf(target).Interface().(*reflect.Value)
	} else {
		v = reflect.ValueOf(target).Elem()
	}
	return v
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
