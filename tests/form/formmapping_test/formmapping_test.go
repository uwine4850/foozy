package formmappingtest

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"slices"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/form/formmapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/tests/common/tconf"
	"github.com/uwine4850/foozy/tests/common/tutils"
)

type TestMapping struct {
	Text []string        `name:"text" empty:"-err"`
	File []form.FormFile `name:"file" empty:"-err"`
}

func mpDefaultStruct(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		panic(err)
	}
	var testMapping TestMapping
	mapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&testMapping), []string{})
	if err := mapper.Fill(); err != nil {
		panic(err)
	}
	if testMapping.Text[0] != "text" {
		w.Write([]byte("the value of the 'text' field does not match the expected value"))
	}
	if testMapping.File[0].Header.Filename != "x.png" {
		w.Write([]byte("the value of the 'file' field does not match the expected value"))
	}
	var testMappingValue TestMapping
	value := reflect.ValueOf(&testMappingValue).Elem()
	valueMapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&value), []string{})
	if err := valueMapper.Fill(); err != nil {
		panic(err)
	}
	if testMappingValue.Text[0] != "text" {
		w.Write([]byte("the value of the 'text' field does not match the expected value"))
	}
	if testMappingValue.File[0].Header.Filename != "x.png" {
		w.Write([]byte("the value of the 'file' field does not match the expected value"))
	}
	return func() {}
}

func TestDefaultForm(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "mp-default-struct"),
		map[string][]string{"text": {"text"}},
		map[string][]string{"file": {"x.png"}},
	)
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}

func mpEmptyString0Err(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		panic(err)
	}
	var testMapping TestMapping
	mapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&testMapping), []string{})
	err := mapper.Fill()
	if !errors.Is(err, formmapper.ErrEmptyFieldIndex{Name: "Text", Index: "0"}) {
		return func() {
			w.Write([]byte("expected error ErrEmptyFieldIndex with field name Text and index 0 not found"))
		}
	}
	var testMappingValue TestMapping
	value := reflect.ValueOf(&testMappingValue).Elem()
	valueMapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&value), []string{})
	err = valueMapper.Fill()
	if !errors.Is(err, formmapper.ErrEmptyFieldIndex{Name: "Text", Index: "0"}) {
		return func() {
			w.Write([]byte("expected error ErrEmptyFieldIndex with field name Text and index 0 not found"))
		}
	}
	return func() {}
}

func TestEmptyString0(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "mp-empty-string-0-err"),
		map[string][]string{"text": {""}},
		map[string][]string{"file": {"x.png"}},
	)
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}

func mpEmptyString1Err(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		panic(err)
	}
	var testMapping TestMapping
	mapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&testMapping), []string{})
	err := mapper.Fill()
	if !errors.Is(err, formmapper.ErrEmptyFieldIndex{Name: "Text", Index: "1"}) {
		return func() {
			w.Write([]byte("expected error ErrEmptyFieldIndex with field name Text and index 1 not found"))
		}
	}
	var testMappingValue TestMapping
	value := reflect.ValueOf(&testMappingValue).Elem()
	valueMapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&value), []string{})
	err = valueMapper.Fill()
	if !errors.Is(err, formmapper.ErrEmptyFieldIndex{Name: "Text", Index: "1"}) {
		return func() {
			w.Write([]byte("expected error ErrEmptyFieldIndex with field name Text and index 1 not found"))
		}
	}
	return func() {}
}

func TestEmptyString1(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "mp-empty-string-1-err"),
		map[string][]string{"text": {"text", ""}},
		map[string][]string{"file": {"x.png"}},
	)
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}

func mpEmptyFileErr(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		panic(err)
	}
	var testMapping TestMapping
	mapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&testMapping), []string{})
	err := mapper.Fill()
	if !errors.Is(err, formmapper.ErrEmptyFieldIndex{Name: "File", Index: "unkown"}) {
		return func() {
			w.Write([]byte("expected error ErrEmptyFieldIndex with field name File and index unkown not found"))
		}
	}
	var testMappingValue TestMapping
	value := reflect.ValueOf(&testMappingValue).Elem()
	valueMapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&value), []string{})
	err = valueMapper.Fill()
	if !errors.Is(err, formmapper.ErrEmptyFieldIndex{Name: "File", Index: "unkown"}) {
		return func() {
			w.Write([]byte("expected error ErrEmptyFieldIndex with field name File and index unkown not found"))
		}
	}
	return func() {}
}

func TestEmptyFile(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "mp-empty-file-err"),
		map[string][]string{"text": {"text"}},
		map[string][]string{"file": {""}},
	)
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}

type TestMappingEmpty struct {
	Text []string `name:"text" empty:"def"`
}

func mpEmptyValue(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		panic(err)
	}
	var testMapping TestMappingEmpty
	mapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&testMapping), []string{})
	if err := mapper.Fill(); err != nil {
		panic(err)
	}
	if testMapping.Text[0] != "def" {
		w.Write([]byte("'empty' The default value is not set."))
	}
	if testMapping.Text[1] != "def" {
		w.Write([]byte("'empty' The default value is not set."))
	}

	var testMappingValue TestMappingEmpty
	value := reflect.ValueOf(&testMappingValue).Elem()
	valueMapper := formmapper.NewMapper(frm, typeopr.Ptr{}.New(&value), []string{})
	if err := mapper.Fill(); err != nil {
		panic(err)
	}
	if err := valueMapper.Fill(); err != nil {
		panic(err)
	}
	if testMappingValue.Text[0] != "def" {
		w.Write([]byte("'empty' The default value is not set."))
	}
	if testMappingValue.Text[1] != "def" {
		w.Write([]byte("'empty' The default value is not set."))
	}
	return func() {}
}

func TestEmptyValue(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "mp-empty-value"),
		map[string][]string{"text": {"", ""}},
		map[string][]string{},
	)
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}

type FillTypedStruct struct {
	Id    int           `name:"id"`
	Name  string        `name:"name"`
	Slice []string      `name:"slice"`
	IsOk  bool          `name:"is_ok"`
	File  form.FormFile `name:"file"`
}

func mpTypedMapper(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	expectData := FillTypedStruct{
		Id:    1,
		Name:  "name",
		Slice: []string{"a1", "a2"},
		IsOk:  true,
	}
	frm := form.NewForm(r)
	if err := frm.Parse(); err != nil {
		panic(err)
	}
	var out FillTypedStruct
	mapper := formmapper.NewTypedMapper(frm, typeopr.Ptr{}.New(&out))
	if err := mapper.Fill(); err != nil {
		panic(err)
	}
	if expectData.Id != out.Id {
		w.Write([]byte("one of the expected data does not match"))
	}
	if expectData.Name != out.Name {
		w.Write([]byte("one of the expected data does not match"))
	}
	if !slices.Equal(expectData.Slice, out.Slice) {
		w.Write([]byte("one of the expected data does not match"))
	}
	if expectData.IsOk != out.IsOk {
		w.Write([]byte("one of the expected data does not match"))
	}
	if out.File.Header.Filename != "x.png" {
		w.Write([]byte("one of the expected data does not match"))
	}
	return func() {}
}

func TestTypedMapper(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortFormMapping, "mp-typed-struct"),
		map[string][]string{"id": {"1"}, "name": {"name"}, "slice": {"a1", "a2"}, "is_ok": {"true"}},
		map[string][]string{"file": {"x.png"}},
	)
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}
