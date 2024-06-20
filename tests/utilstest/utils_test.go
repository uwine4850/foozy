package utilstest

import (
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces/intrnew"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type ITStruct interface {
	intrnew.INewInstance
	Method()
}

type TStruct struct{}

func (t *TStruct) New() (interface{}, error) {
	return &TStruct{}, nil
}

func (t *TStruct) Method() {}

func TestCreateNewInstance(t *testing.T) {
	ff := &TStruct{}
	var newStruct ITStruct
	if err := typeopr.CreateNewInstance(ff, &newStruct); err != nil {
		t.Error(err)
	}
}
