package objecttest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/object"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
)

type JsonObjectViewMessage struct {
	rest.ImplementDTOMessage
	Id   string `json:"Id"`
	Name string `json:"Name"`
	FF   string `json:"Ff"`
	Test int    `json:"Test"`
}

type JsonFormMessage struct {
	rest.ImplementDTOMessage
	Text string `json:"Text"`
	Test string `json:"Test"`
}

var dto = rest.NewDTO()

func TestMain(m *testing.M) {
	dto.AllowedMessages([]rest.AllowMessage{
		{
			Package: "objecttest",
			Name:    "JsonObjectViewMessage",
		},
		{
			Package: "objecttest",
			Name:    "JsonFormMessage",
		},
	})
	dto.Messages(map[string]*[]irest.IMessage{
		"s": {
			JsonObjectViewMessage{},
			JsonFormMessage{},
		},
	})

	db := database.NewDatabase(database.DbArgs{Username: "root", Password: "1111", Host: "localhost", Port: "3408", DatabaseName: "foozy_test"})
	if err := db.Connect(); err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	createAndFillTable(db)

	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newRouter := router.NewRouter(manager.NewManager(render))
	newRouter.Get("/object-view/<id>", TObjectViewHNDL(db))
	newRouter.Get("/object-mul-view/<id>/<id1>", TObjectMultipleViewHNDL(db))
	newRouter.Get("/object-all-view", TObjectAllViewHNDL(db))
	newRouter.Post("/object-form-view", MyFormViewHNDL())
	newRouter.Get("/object-json-view/<id>", TJsonObjectViewHNDL(db))
	newRouter.Get("/object-mul-json-view/<id>/<id1>", TJsonObjectMultipleViewHNDL(db))
	newRouter.Get("/object-json-all-view", TJsonObjectAllViewHNDL(db))

	serv := server.NewServer(":8031", newRouter, nil)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(":8031", 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

var tableQuery = "CREATE TABLE IF NOT EXISTS object_test (" +
	"`id` INT NOT NULL AUTO_INCREMENT, " +
	"`name` VARCHAR(200) NOT NULL, " +
	"`FF` FLOAT NULL, " +
	"PRIMARY KEY (id)" +
	");"

var tableQuery1 = "CREATE TABLE IF NOT EXISTS object_test1 (" +
	"`id` INT NOT NULL AUTO_INCREMENT, " +
	"`name` VARCHAR(200) NOT NULL, " +
	"`FF` FLOAT NULL, " +
	"PRIMARY KEY (id)" +
	");"

func createAndFillTable(db *database.Database) {
	if _, err := db.SyncQ().Query(tableQuery); err != nil {
		panic(err)
	}
	if _, err := db.SyncQ().Query("TRUNCATE TABLE object_test;"); err != nil {
		panic(err)
	}
	if _, err := db.SyncQ().Insert("object_test", map[string]interface{}{"name": "name"}); err != nil {
		panic(err)
	}
	if _, err := db.SyncQ().Insert("object_test", map[string]interface{}{"name": "name0"}); err != nil {
		panic(err)
	}

	if _, err := db.SyncQ().Query(tableQuery1); err != nil {
		panic(err)
	}
	if _, err := db.SyncQ().Query("TRUNCATE TABLE object_test1;"); err != nil {
		panic(err)
	}
	if _, err := db.SyncQ().Insert("object_test1", map[string]interface{}{"name": "name1"}); err != nil {
		panic(err)
	}
}

type TObjectViewDB struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	FF   string `db:"FF"`
}

type TObjectView struct {
	object.ObjView
}

func (v *TObjectView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *TObjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	objectContext, err := object.GetObjectContext(manager)
	if err != nil {
		panic(err)
	}
	if _, ok := objectContext["object"]; !ok {
		panic("ObjectContext does not have a key object.")
	}
	return map[string]interface{}{"TEST": "OK"}, nil
}

func TObjectViewHNDL(db *database.Database) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.TemplateView{
		TemplatePath: "./templates/object_view.html",
		View: &TObjectView{
			object.ObjView{
				Name:       "object",
				DB:         db,
				TableName:  "object_test",
				FillStruct: TObjectViewDB{},
				Slug:       "id",
			},
		},
	}
	return view.Call
}

func TestObjectView(t *testing.T) {
	get, err := http.Get("http://localhost:8031/object-view/1")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "1 OK" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

type TObjectMultipleView struct {
	object.MultipleObjectView
}

func (v *TObjectMultipleView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *TObjectMultipleView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["object"]; !ok {
		panic("ObjectContext does not have a key object.")
	}
	if _, ok := objectContext["object1"]; !ok {
		panic("ObjectContext does not have a key object1.")
	}
	return object.ObjectContext{"TEST": "OK"}, nil
}

func TObjectMultipleViewHNDL(db *database.Database) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.TemplateView{
		TemplatePath: "./templates/object_multiple_view.html",
		View: &TObjectMultipleView{
			object.MultipleObjectView{
				DB: db,
				MultipleObjects: []object.MultipleObject{
					{
						Name:       "object",
						TaleName:   "object_test",
						SlugName:   "id",
						SlugField:  "id",
						FillStruct: TObjectViewDB{},
					},
					{
						Name:       "object1",
						TaleName:   "object_test1",
						SlugName:   "id1",
						SlugField:  "id",
						FillStruct: TObjectViewDB{},
					},
				},
			},
		},
	}
	return view.Call
}

func TestObjectMultipleView(t *testing.T) {
	get, err := http.Get("http://localhost:8031/object-mul-view/1/1")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "name name1 OK" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

type TObjectAllView struct {
	object.AllView
}

func (v *TObjectAllView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *TObjectAllView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["all_object"]; !ok {
		panic("ObjectContext does not have a key all_object.")
	}
	return object.ObjectContext{"TEST": "OK"}, nil
}

func TObjectAllViewHNDL(db *database.Database) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.TemplateView{
		TemplatePath: "./templates/object_all_view.html",
		View: &TObjectAllView{
			object.AllView{
				Name:       "all_object",
				DB:         db,
				TableName:  "object_test",
				FillStruct: TObjectViewDB{},
			},
		},
	}
	return view.Call
}

func TestObjectAllView(t *testing.T) {
	get, err := http.Get("http://localhost:8031/object-all-view")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "name name0 OK" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

type ObjectForm struct {
	Text []string        `form:"text"`
	File []form.FormFile `form:"file" ext:".jpg .png"`
}

type MyFormView struct {
	object.FormView
}

func (v *MyFormView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	filledFormInterface, err := v.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	filledForm := filledFormInterface.(ObjectForm)
	if filledForm.Text[0] != "field" {
		return nil, errors.New("FormView unexpected text field value")
	}
	if filledForm.File[0].Header.Filename != "x.png" {
		return nil, errors.New("FormView unexpected file field value")
	}
	return object.ObjectContext{}, nil
}

func (v *MyFormView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func MyFormViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	tv := object.TemplateView{
		TemplatePath: "",
		View: &MyFormView{
			object.FormView{
				FormStruct:       ObjectForm{},
				NotNilFormFields: []string{"*"},
				NilIfNotExist:    []string{},
			},
		},
	}
	tv.SkipRender()
	return tv.Call
}

func TestMyFormView(t *testing.T) {
	multipartForm, err := form.SendMultipartForm("http://localhost:8031/object-form-view", map[string][]string{"text": {"field"}}, map[string][]string{"file": {"x.png"}})
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

type JsonObjectView struct {
	object.ObjView
}

func (v *JsonObjectView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *JsonObjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	objectContext, err := object.GetObjectContext(manager)
	if err != nil {
		panic(err)
	}
	if _, ok := objectContext["object"]; !ok {
		panic("ObjectContext does not have a key object.")
	}
	return map[string]interface{}{"Test": 1}, nil
}

func TJsonObjectViewHNDL(db *database.Database) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.JsonObjectTemplateView{
		View: &TObjectView{
			object.ObjView{
				Name:       "object",
				DB:         db,
				TableName:  "object_test",
				FillStruct: TObjectViewDB{},
				Slug:       "id",
			},
		},
		DTO:     dto,
		Message: JsonObjectViewMessage{},
	}
	return view.Call
}

func TestJsonObjectView(t *testing.T) {
	get, err := http.Get("http://localhost:8031/object-json-view/1")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"Id":"1","Name":"name","Ff":"","Test":0}` {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

type TJsonObjectMultipleView struct {
	object.MultipleObjectView
}

func (v *TJsonObjectMultipleView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *TJsonObjectMultipleView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["object"]; !ok {
		panic("ObjectContext does not have a key object.")
	}
	if _, ok := objectContext["object1"]; !ok {
		panic("ObjectContext does not have a key object1.")
	}
	return object.ObjectContext{"Test": "OK"}, nil
}

func TJsonObjectMultipleViewHNDL(db *database.Database) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.JsonMultipleObjectTemplateView{
		View: &TObjectMultipleView{
			object.MultipleObjectView{
				DB: db,
				MultipleObjects: []object.MultipleObject{
					{
						Name:       "object",
						TaleName:   "object_test",
						SlugName:   "id",
						SlugField:  "id",
						FillStruct: TObjectViewDB{},
					},
					{
						Name:       "object1",
						TaleName:   "object_test1",
						SlugName:   "id1",
						SlugField:  "id",
						FillStruct: TObjectViewDB{},
					},
				},
			},
		},
		DTO:     dto,
		Message: JsonObjectViewMessage{},
	}
	return view.Call
}

func TestJsonObjectMultipleView(t *testing.T) {
	get, err := http.Get("http://localhost:8031/object-mul-json-view/1/1")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `[{"Id":"1","Name":"name","Ff":"","Test":0},{"Id":"1","Name":"name1","Ff":"","Test":0}]` {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

type TJsonObjectAllView struct {
	object.AllView
}

func (v *TJsonObjectAllView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *TJsonObjectAllView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["all_object"]; !ok {
		panic("ObjectContext does not have a key all_object.")
	}
	return object.ObjectContext{"Test": "OK"}, nil
}

func TJsonObjectAllViewHNDL(db *database.Database) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.JsonAllTemplateView{
		View: &TObjectAllView{
			object.AllView{
				Name:       "all_object",
				DB:         db,
				TableName:  "object_test",
				FillStruct: TObjectViewDB{},
			},
		},
		DTO:     dto,
		Message: JsonObjectViewMessage{},
	}
	return view.Call
}

func TestJsonObjectAllView(t *testing.T) {
	get, err := http.Get("http://localhost:8031/object-json-all-view")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(body))
	// if string(body) != "name name0 OK" {
	// 	t.Errorf("Error on page retrieval.")
	// }
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
