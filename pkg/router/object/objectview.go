package object

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
	"reflect"
)

type ErrNoSlug struct {
	SlugName string
}

func (e ErrNoSlug) Error() string {
	return fmt.Sprintf("slug parameter %s not found", e.SlugName)
}

type ErrNoData struct {
}

func (e ErrNoData) Error() string {
	return "no data to display was found"
}

// ObjView displays only the HTML page only with a specific row from the database.
// Needs to be used with slug parameter URL path, specify the name of the parameter in the Slug parameter.
type ObjView struct {
	IView

	Name       string
	DB         *database.Database
	TableName  string
	FillStruct interface{}
	Slug       string

	context map[string]interface{}
}

func (v *ObjView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *ObjView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
	return map[string]interface{}{}
}

func (v *ObjView) GetContext() map[string]interface{} {
	return v.context
}

func (v *ObjView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error) {
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
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
	if res == nil {
		return nil, ErrNoData{}
	}
	value, err := v.fillObject(res[0])
	if err != nil {
		return nil, err
	}
	v.context = map[string]interface{}{v.Name: value.Interface()}
	// CLOSE DB
	err = v.DB.Close()
	if err != nil {
		return nil, err
	}
	return v.context, nil
}

func (v *ObjView) fillObject(object map[string]interface{}) (*reflect.Value, error) {
	if v.FillStruct == nil {
		panic("the FillStruct field must not be nil")
	}
	value := reflect.New(reflect.TypeOf(v.FillStruct)).Elem()
	err := dbutils.FillReflectValueFromDb(object, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (v *ObjView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic("OnError is not implement")
}
