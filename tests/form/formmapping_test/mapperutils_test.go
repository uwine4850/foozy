package formmappingtest

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/form/formmapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/tests1/common/tconf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

type Fill struct {
	NilField []string `form:"isNil"`
	Str      string
	Field1   []string        `form:"f1"`
	File     []form.FormFile `form:"file" empty:""`
}

func fill(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	var f Fill
	err = formmapper.FillStructFromForm(newForm, typeopr.Ptr{}.New(&f), []string{"isNil"})
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	if f.NilField != nil {
		return func() { w.Write([]byte("The NilField must be nil.")) }
	}
	if f.Field1 == nil {
		return func() { w.Write([]byte("The Field1 field must be populated.")) }
	}
	if f.File == nil {
		return func() { w.Write([]byte("The File field must be populated.")) }
	}
	return func() {}
}

func TestFillStructFromForm(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "fill"),
		map[string][]string{"f1": {"v1"}}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
	err = multipartForm.Body.Close()
	if err != nil {
		panic(err)
	}
}

func fillReflectValue(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	var f Fill
	fValue := reflect.ValueOf(&f).Elem()
	err = formmapper.FillReflectValueFromForm(newForm, &fValue, []string{"isNil"})
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	if f.NilField != nil {
		return func() { w.Write([]byte("The NilField must be nil.")) }
	}
	if f.Field1 == nil {
		return func() { w.Write([]byte("The Field1 field must be populated.")) }
	}
	if f.File == nil {
		return func() { w.Write([]byte("The File field must be populated.")) }
	}
	return func() {}
}

func TestFillReflectValueFromForm(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "fill-reflect-value"),
		map[string][]string{"f1": {"v1"}}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
	err = multipartForm.Body.Close()
	if err != nil {
		panic(err)
	}
}
