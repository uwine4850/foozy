package typeopr_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

func TestIsPointer(t *testing.T) {
	value := "val"
	if !typeopr.IsPointer(&value) {
		t.Error("the value is actually a pointer")
	}
}

func TestPtrIsStruct(t *testing.T) {
	s := struct{}{}
	if !typeopr.PtrIsStruct(&s) {
		t.Error("the value is actually a pointer to a structure")
	}
	value := "val"
	if typeopr.PtrIsStruct(&value) {
		fmt.Println("the value is not actually a pointer to a structure")
	}
}

func TestIsEmpty(t *testing.T) {
	empty := []string{}
	if !typeopr.IsEmpty(empty) {
		t.Error("is actually an empty value")
	}
}

func TestAnyToBytes(t *testing.T) {
	value := "VALUE"
	byteRes, err := typeopr.AnyToBytes(value)
	if err != nil {
		t.Error(err)
	}
	if string(byteRes) != "VALUE" {
		t.Error("byte conversion error")
	}
}

func TestPtrObject(t *testing.T) {
	value := 1
	pointer := typeopr.Ptr{}.New(&value)
	if reflect.TypeOf(pointer.Ptr()).Kind() != reflect.Pointer {
		t.Error("value must be a pointer")
	}
}

type FakePointer struct{}

func (f FakePointer) New(val interface{}) typeopr.IPtr {
	return f
}

func (f FakePointer) Ptr() interface{} {
	return ""
}

func TestIsImplementInterface(t *testing.T) {
	object := FakePointer{}
	if !typeopr.IsImplementInterface(typeopr.Ptr{}.New(&object), (*typeopr.IPtr)(nil)) {
		t.Error("the object actually implements the interface")
	}
	nilObject := struct{}{}
	if typeopr.IsImplementInterface(typeopr.Ptr{}.New(&nilObject), (*typeopr.IPtr)(nil)) {
		t.Error("the check is not correct, the interface is not implemented")
	}
}

func TestGetReflectValue(t *testing.T) {
	val := 111
	reflectValue := typeopr.GetReflectValue(&val)
	if !reflect.DeepEqual(reflect.TypeOf(reflectValue), reflect.TypeOf(reflect.Value{})) {
		t.Error("conversion to reflect.Value is not correct")
	}
}
