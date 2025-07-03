package utilst_test

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"testing"

	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
	"github.com/uwine4850/foozy/pkg/utils/fstruct"
)

func TestMergeMap(t *testing.T) {
	m1 := map[string]string{"1": "OK1", "2": "OK2"}
	m2 := map[string]string{"11": "OK11", "22": "OK2"}
	fmap.MergeMap(&m1, m2)
	expectedMap := map[string]string{"1": "OK1", "2": "OK2", "11": "OK11", "22": "OK2"}
	if !maps.Equal(expectedMap, m1) {
		t.Error("the expected map and the passed map do not match")
	}
}

type YamlObject struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
	Ok   bool   `yaml:"ok"`
}

func TestYamlMapToStruct(t *testing.T) {
	yamlMap := map[string]any{"id": 1, "name": "TestName", "ok": true}
	var yamlObject YamlObject
	if err := fmap.YamlMapToStruct(&yamlMap, typeopr.Ptr{}.New(&yamlObject)); err != nil {
		t.Error(err)
	}
	expectedYamlObject := YamlObject{
		Id:   1,
		Name: "TestName",
		Ok:   true,
	}
	if !reflect.DeepEqual(yamlObject, expectedYamlObject) {
		t.Error("the filled and expected structure do not match")
	}
}

func TestPathExists(t *testing.T) {
	if !fpath.PathExist("./dir/file.txt") {
		t.Error("path is valid but is displayed as not found")
	}
}

func TestSliceContains(t *testing.T) {
	sl := []int{1, 2, 3}
	if !fslice.SliceContains(sl, 1) {
		t.Error("slice contains the value")
	}
	if fslice.SliceContains(sl, 111) {
		t.Error("slice does not contain a value")
	}
}

func TestAllStringItemsEmpty(t *testing.T) {
	sl := []string{"", "", ""}
	if !fslice.AllStringItemsEmpty(sl) {
		t.Error("is actually a slice with empty values")
	}
	sl2 := []string{"11", "1", ""}
	if fslice.AllStringItemsEmpty(sl2) {
		t.Error("is not actually a slice with empty values.")
	}
}

func TestToLower(t *testing.T) {
	text := "TEXT"
	if fstring.ToLower(text) != "tEXT" {
		t.Error("ToLower function didn't work")
	}
}

type SomeStruct struct {
	Id   int `notdef:"true"`
	Name string
}

func TestCheckNotDefaultFields(t *testing.T) {
	someStruct1 := SomeStruct{
		Id:   1,
		Name: "NAME",
	}
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&someStruct1)); err != nil {
		t.Error(err)
	}
	someStruct2 := SomeStruct{
		Name: "NAME",
	}
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&someStruct2)); err != nil {
		if !errors.As(err, &fstruct.ErrStructFieldIsDefault{}) {
			fmt.Println("expected error not detected")
		}
	}
}
