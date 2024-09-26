package utilstest

import (
	"errors"
	"testing"

	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fstruct"
)

type Tee struct {
	Id      int            `notdef:"true"`
	Name    string         `notdef:"true"`
	Slice   []int          `notdef:"true"`
	Map     map[string]int `notdef:"true"`
	Pointer *int           `notdef:"true"`
}

func TestCheckNotDefaultFields(t *testing.T) {
	p := 1
	tee := Tee{
		Id:      1,
		Name:    "name",
		Slice:   []int{1, 2, 3},
		Map:     map[string]int{"1": 1, "2": 2, "3": 3},
		Pointer: &p,
	}
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&tee)); err != nil {
		t.Error(err)
	}
}

func TestCheckNotDefaultFieldsError(t *testing.T) {
	p := 1
	tee := Tee{
		Id:      1,
		Slice:   []int{1, 2, 3},
		Map:     map[string]int{"1": 1, "2": 2, "3": 3},
		Pointer: &p,
	}
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&tee)); err != nil {
		if !errors.Is(err, fstruct.ErrStructFieldIsDefault{FieldName: "Name"}) {
			t.Error(err)
		}
	} else {
		t.Error("CheckNotDefaultFields should return an error.")
	}
	tee1 := Tee{
		Id:      1,
		Name:    "name",
		Slice:   []int{1, 2, 3},
		Pointer: &p,
	}
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&tee1)); err != nil {
		if !errors.Is(err, fstruct.ErrStructFieldIsDefault{FieldName: "Map"}) {
			t.Error(err)
		}
	} else {
		t.Error("CheckNotDefaultFields should return an error.")
	}
}
