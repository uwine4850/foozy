package object

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
	"reflect"
)

type ErrNoSlug struct {
	SlugName string
}

func (e ErrNoSlug) Error() string {
	return fmt.Sprintf("slug parameter %s not found", e.SlugName)
}

// ObjView displays only the HTML page only with a specific row from the database.
// Needs to be used with slug parameter URL path, specify the name of the parameter in the Slug parameter.
type ObjView struct {
	UserView
	Name         string
	TemplatePath string
	DB           *database.Database
	TableName    string
	FillStruct   interface{}
	Slug         string

	onError func(err error)
}

func (v *ObjView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *ObjView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
	return map[string]interface{}{}
}

func (v *ObjView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error) {
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			v.onError(err)
		}
	}(v.DB)
	slugValue, ok := manager.GetSlugParams(v.Slug)
	if !ok {
		return nil, ErrNoSlug{v.Slug}
	}
	res, err := v.DB.SyncQ().Select([]string{"*"}, v.TableName, dbutils.WHEquals(map[string]interface{}{
		v.Slug: slugValue,
	}, "AND"), 1)
	if err != nil {
		return nil, err
	}
	object, err := v.fillObject(res[0])
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{v.Name: object}, nil
}

func (v *ObjView) fillObject(object map[string]interface{}) (TemplateStruct, error) {
	if v.FillStruct != nil {
		value := reflect.New(reflect.TypeOf(v.FillStruct)).Elem()
		err := dbutils.FillReflectValueFromDb(object, &value)
		if err != nil {
			return TemplateStruct{}, err
		}
		return TemplateStruct{s: value}, nil
	}
	fillMap := make(map[string]string)
	err := dbutils.FillMapFromDb(object, &fillMap)
	if err != nil {
		return TemplateStruct{}, err
	}
	return TemplateStruct{m: fillMap}, nil
}

func (v *ObjView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.UserView == nil {
		panic("the UserView field must not be nil")
	}
	permissions, f := v.UserView.Permissions(w, r, manager)
	if !permissions {
		return func() { f() }
	}
	context := v.UserView.Context(w, r, manager)
	object, err := v.Object(w, r, manager)
	if err != nil {
		return func() { v.onError(err) }
	}
	utils.MergeMap(&context, object)
	manager.SetContext(context)
	manager.SetTemplatePath(v.TemplatePath)
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { v.onError(err) }
	}
	return func() {}
}

func (v *ObjView) OnError(e func(err error)) {
	v.onError = e
}
