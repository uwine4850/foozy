package object_test_1

import (
	"database/sql/driver"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/router/object"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

type DTOObjectView struct {
	object.ObjView
}

func (v *DTOObjectView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func jsonObjectView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectWhereMock(mock, "table", "1", [][]driver.Value{
		{1, "TEST_NAME", true},
	})
	db := NewMockDatabase(sqlDB)
	view := object.JsonObjectTemplateView{
		View: &DTOObjectView{
			object.ObjView{
				Name:       "object",
				TableName:  "table",
				Database:   db,
				FillStruct: DatabaseTable{},
				Slug:       "id",
			},
		},
		DTO:     newDTO,
		Message: DTOMessage{},
	}
	return view.Call
}

func TestJsonObjectView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-json-object-view/1"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != `{"TypDTOMessage":{},"Id":1,"Name":"TEST_NAME","Ok":true}` {
		t.Error("json object view error: ", res)
	}
}

type JsonAllView struct {
	object.AllView
}

func (v *JsonAllView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectJsonAllView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectAllMock(mock)
	db := NewMockDatabase(sqlDB)
	view := object.JsonAllTemplateView{
		View: &JsonAllView{
			object.AllView{
				Name:       "objects",
				TableName:  "table",
				Database:   db,
				FillStruct: DatabaseTable{},
			},
		},
		DTO:     newDTO,
		Message: DTOMessage{},
	}
	return view.Call
}

func TestJsonAllView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-json-all-view"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != `[{"TypDTOMessage":{},"Id":1,"Name":"TEST_NAME","Ok":true},{"TypDTOMessage":{},"Id":2,"Name":"TEST_NAME_1","Ok":true}]` {
		t.Error("object view error: ", res)
	}
}

type JsonSlugAllView struct {
	object.AllView
}

func (v *JsonSlugAllView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectSlugJsonAllView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	initSelectWhereMock(mock, "table", "1", [][]driver.Value{
		{1, "TEST_NAME", true},
	})
	db := NewMockDatabase(sqlDB)
	view := object.JsonAllTemplateView{
		View: &JsonAllView{
			object.AllView{
				Name:       "objects",
				TableName:  "table",
				Database:   db,
				Slug:       "id",
				FillStruct: DatabaseTable{},
			},
		},
		DTO:     newDTO,
		Message: DTOMessage{},
	}
	return view.Call
}

func TestSlugJsonAllView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-slug-json-all-view/1"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != `[{"TypDTOMessage":{},"Id":1,"Name":"TEST_NAME","Ok":true}]` {
		t.Error("object view error: ", res)
	}
}

type JsonMultipleView struct {
	object.MultipleObjectView
}

func (v *JsonMultipleView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func objectJsonMultipleView() func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
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
	view := object.JsonMultipleObjectTemplateView{
		View: &JsonMultipleView{
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
		DTO: newDTO,
		Messages: map[string]irest.Message{
			"object1": DTOMessage{},
			"object2": DTOMessage{},
		},
	}
	return view.Call
}

func TestJsomMultipleView(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortObject, "test-json-multiple-view/obj1/obj2"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != `{"object1":{"TypDTOMessage":{},"Id":1,"Name":"TEST_NAME","Ok":true},"object2":{"TypDTOMessage":{},"Id":2,"Name":"TEST_NAME","Ok":true}}` {
		t.Error("object view error: ", res)
	}
}
