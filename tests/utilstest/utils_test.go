package utilstest

import (
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
)

type ITStruct interface {
	interfaces.INewInstance
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
	if err := utils.CreateNewInstance(ff, &newStruct); err != nil {
		t.Error(err)
	}
}
