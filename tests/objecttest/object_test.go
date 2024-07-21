package objecttest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/object"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
)

func TestMain(m *testing.M) {
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
	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.Get("/object-view/<id>", TObjectViewHNDL(db))
	newRouter.Get("/object-mul-view/<id>/<id1>", TObjectMultipleViewHNDL(db))
	newRouter.Get("/object-all-view", TObjectAllViewHNDL(db))
	newRouter.Post("/object-form-view", MyFormViewHNDL())

	serv := server.NewServer(":8030", newRouter)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(":8030", 5); err != nil {
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
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["object"]; !ok {
		panic("the context has object data")
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
	get, err := http.Get("http://localhost:8030/object-view/1")
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
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["object"]; !ok {
		panic("the context has object data")
	}
	if _, ok := objectContext["object1"]; !ok {
		panic("the context has object data")
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
	get, err := http.Get("http://localhost:8030/object-mul-view/1/1")
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
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.ObjectContext)
	if _, ok := objectContext["all_object"]; !ok {
		panic("the context has object data")
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
	get, err := http.Get("http://localhost:8030/object-all-view")
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
	multipartForm, err := form.SendMultipartForm("http://localhost:8030/object-form-view", map[string]string{"text": "field"}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Errorf(string(responseBody))
	}
}
