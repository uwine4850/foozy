package object_test_1

import (
	"database/sql/driver"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/object"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

// OBJECT VIEW TEST -------------------------------------

type ObjectView struct {
	object.ObjView
}

func (v *ObjectView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectWhereMock(mock, "table", "1", [][]driver.Value{
		{1, "TEST_NAME", true},
	})
	db := NewMockDatabase(sqlDB)
	view := object.TemplateView{
		TemplatePath: "object_template.html",
		View: &ObjectView{
			object.ObjView{
				Name:       "object",
				TableName:  "table",
				Database:   db,
				FillStruct: DatabaseTable{},
				Slug:       "id",
			},
		},
	}
	return view.Call
}

func TestObjectView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-object-view/1"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "1" {
		t.Error("object view error: ", res)
	}
}

// ALL VIEW TEST ---------------------

type AllView struct {
	object.AllView
}

func (v *AllView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectAllView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectAllMock(mock)
	db := NewMockDatabase(sqlDB)
	view := object.TemplateView{
		TemplatePath: "all_template.html",
		View: &AllView{
			object.AllView{
				Name:       "objects",
				TableName:  "table",
				Database:   db,
				FillStruct: DatabaseTable{},
			},
		},
	}
	return view.Call
}

func TestAllView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-all-view"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "1|2|" {
		t.Error("object view error: ", res)
	}
}

// ALL VIEW WITH SLUG TEST ----------------------

type AllSlugView struct {
	object.AllView
}

func (v *AllSlugView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectAllSlugView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectWhereMock(mock, "table", "1", [][]driver.Value{
		{1, "TEST_NAME", true},
	})
	db := NewMockDatabase(sqlDB)
	view := object.TemplateView{
		TemplatePath: "all_template.html",
		View: &AllSlugView{
			object.AllView{
				Name:       "objects",
				TableName:  "table",
				Database:   db,
				FillStruct: DatabaseTable{},
				Slug:       "id",
			},
		},
	}
	return view.Call
}

func TestAllSlugView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-all-view/1"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "1|" {
		t.Error("object view error: ", res)
	}
}

// MULTIPLE VIEW TEST ----------------------

type MultipleView struct {
	object.MultipleObjectView
}

func (v *MultipleView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectMultipleView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectWhereMock(mock, "table1", "obj1", [][]driver.Value{
		{1, "TEST_NAME", true},
	})
	initSelectWhereMock(mock, "table2", "obj2", [][]driver.Value{
		{2, "TEST_NAME", true},
	})
	db := NewMockDatabase(sqlDB)
	view := object.TemplateView{
		TemplatePath: "multiple_template.html",
		View: &MultipleView{
			object.MultipleObjectView{
				Database: db,
				MultipleObjects: []object.MultipleObject{
					{
						Name:       "object1",
						TableName:  "table1",
						SlugName:   "o1",
						SlugField:  "id",
						FillStruct: DatabaseTable{},
					},
					{
						Name:       "object2",
						TableName:  "table2",
						SlugName:   "o2",
						SlugField:  "id",
						FillStruct: DatabaseTable{},
					},
				},
			},
		},
	}
	return view.Call
}

func TestMultipleView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-multiple-view/obj1/obj2"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "1|2" {
		t.Error("object view error: ", res)
	}
}
