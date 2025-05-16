package objecttest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
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
	"github.com/uwine4850/foozy/tests/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests/common/tutils"
)

type JsonObjectViewMessage struct {
	rest.ImplementDTOMessage
	Id   int     `json:"Id" dto:"Id"`
	Name string  `json:"Name" dto:"Name"`
	FF   float64 `json:"Ff" dto:"Ff"`
	Test string  `json:"Test" dto:"Test"`
}

type JsonFormMessage struct {
	rest.ImplementDTOMessage
	Text string `json:"Text" dto:"Text"`
	Test string `json:"Test" dto:"Test"`
}

var dto = rest.NewDTO()

func TestMain(m *testing.M) {
	testinitcnf.InitCnf()
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
	dto.Messages(map[string][]irest.IMessage{
		"s": {
			JsonObjectViewMessage{},
			JsonFormMessage{},
		},
	})

	db := database.NewDatabase(tconf.DbArgs)
	if err := db.Open(); err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newManager := manager.NewManager(render)
	if err := database.InitDatabasePool(newManager, db); err != nil {
		panic(err)
	}
	dbRead, err := newManager.Database().ConnectionPool(config.LoadedConfig().Default.Database.MainConnectionPoolName)
	if err != nil {
		panic(err)
	}
	createAndFillTable(dbRead)

	newRouter := router.NewRouter(newManager)
	newRouter.Get("/object-view/<id>", TObjectViewHNDL())
	newRouter.Get("/object-mul-view/<id>/<id1>", TObjectMultipleViewHNDL())
	newRouter.Get("/object-all-view", TObjectAllViewHNDL())
	newRouter.Post("/object-form-view", MyFormViewHNDL())
	newRouter.Get("/object-json-view/<id>", TJsonObjectViewHNDL())
	newRouter.Get("/object-mul-json-view/<id>/<id1>", TJsonObjectMultipleViewHNDL())
	newRouter.Get("/object-json-all-view", TJsonObjectAllViewHNDL())

	serv := server.NewServer(tconf.PortObject, newRouter, nil)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortObject, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
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

func createAndFillTable(dbRead interfaces.IReadDatabase) {
	if _, err := dbRead.SyncQ().Query(tableQuery); err != nil {
		panic(err)
	}
	if _, err := dbRead.SyncQ().Query("TRUNCATE TABLE object_test;"); err != nil {
		panic(err)
	}
	q := qb.NewSyncQB(dbRead.SyncQ()).Insert("object_test", map[string]interface{}{"name": "name"})
	q.Merge()
	if _, err := q.Exec(); err != nil {
		panic(err)
	}
	q1 := qb.NewSyncQB(dbRead.SyncQ()).Insert("object_test", map[string]interface{}{"name": "name0"})
	q1.Merge()
	if _, err := q1.Exec(); err != nil {
		panic(err)
	}

	if _, err := dbRead.SyncQ().Query(tableQuery1); err != nil {
		panic(err)
	}
	if _, err := dbRead.SyncQ().Query("TRUNCATE TABLE object_test1;"); err != nil {
		panic(err)
	}
	q2 := qb.NewSyncQB(dbRead.SyncQ()).Insert("object_test1", map[string]interface{}{"name": "name1"})
	q2.Merge()
	if _, err := q2.Exec(); err != nil {
		panic(err)
	}
}

type TObjectViewDB struct {
	Id   int     `db:"id"`
	Name string  `db:"name"`
	FF   float64 `db:"FF"`
}

type TObjectView struct {
	object.ObjView
}

func (v *TObjectView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func (v *TObjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	objectContext, err := object.GetContext(manager)
	if err != nil {
		panic(err)
	}
	if _, ok := objectContext["object"]; !ok {
		panic("ObjectContext does not have a key object.")
	}
	return map[string]interface{}{"TEST": "OK"}, nil
}

func TObjectViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.TemplateView{
		TemplatePath: "./templates/object_view.html",
		View: &TObjectView{
			object.ObjView{
				Name:       "object",
				TableName:  "object_test",
				FillStruct: TObjectViewDB{},
				Slug:       "id",
			},
		},
	}
	return view.Call
}

func TestObjectView(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortObject, "object-view/1"))
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

func (v *TObjectMultipleView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.Context)
	if _, ok := objectContext["object"]; !ok {
		panic("ObjectContext does not have a key object.")
	}
	if _, ok := objectContext["object1"]; !ok {
		panic("ObjectContext does not have a key object1.")
	}
	return object.Context{"TEST": "OK"}, nil
}

func TObjectMultipleViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.TemplateView{
		TemplatePath: "./templates/object_multiple_view.html",
		View: &TObjectMultipleView{
			object.MultipleObjectView{
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortObject, "object-mul-view/1/1"))
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

func (v *TObjectAllView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	_objectContext, _ := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	objectContext := _objectContext.(object.Context)
	if _, ok := objectContext["all_object"]; !ok {
		panic("ObjectContext does not have a key all_object.")
	}
	return object.Context{"TEST": "OK"}, nil
}

func TObjectAllViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.TemplateView{
		TemplatePath: "./templates/object_all_view.html",
		View: &TObjectAllView{
			object.AllView{
				Name:       "all_object",
				TableName:  "object_test",
				FillStruct: TObjectViewDB{},
			},
		},
	}
	return view.Call
}

func TestObjectAllView(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortObject, "object-all-view"))
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

func (v *MyFormView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
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
	return object.Context{}, nil
}

func (v *MyFormView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func MyFormViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	tv := object.TemplateView{
		TemplatePath: "",
		View: &MyFormView{
			object.FormView{
				FormStruct: ObjectForm{},
				// NotNilFormFields: []string{"*"},
				// NilIfNotExist:    []string{},
			},
		},
	}
	tv.SkipRender()
	return tv.Call
}

func TestMyFormView(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortObject, "object-form-view"), map[string][]string{"text": {"field"}}, map[string][]string{"file": {"x.png"}})
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

func (v *JsonObjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	return map[string]interface{}{"Test": "OK"}, nil
}

func TJsonObjectViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.JsonObjectTemplateView{
		View: &JsonObjectView{
			object.ObjView{
				Name:       "object",
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortObject, "object-json-view/1"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"Id":1,"Name":"name","Ff":0,"Test":"OK"}` {
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

func (v *TJsonObjectMultipleView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	return object.Context{"Test": "OK"}, nil
}

func TJsonObjectMultipleViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.JsonMultipleObjectTemplateView{
		View: &TJsonObjectMultipleView{
			object.MultipleObjectView{
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
		DTO: dto,
		Messages: map[string]irest.IMessage{
			"object":  JsonObjectViewMessage{},
			"object1": JsonObjectViewMessage{},
		},
	}
	return view.Call
}

func TestJsonObjectMultipleView(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortObject, "object-mul-json-view/1/1"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"object":{"Id":1,"Name":"name","Ff":0,"Test":""},"object1":{"Id":1,"Name":"name1","Ff":0,"Test":""}}` {
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

func (v *TJsonObjectAllView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	return object.Context{"Test": "OK"}, nil
}

func TJsonObjectAllViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.JsonAllTemplateView{
		View: &TJsonObjectAllView{
			object.AllView{
				Name:       "all_object",
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortObject, "object-json-all-view"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `[{"Id":1,"Name":"name","Ff":0,"Test":"OK"},{"Id":2,"Name":"name0","Ff":0,"Test":"OK"}]` {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
