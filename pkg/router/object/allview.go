package object

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
	"reflect"
)

// AllView displays HTML page by passing all data from the selected table to it.
type AllView struct {
	IView
	Name       string
	DB         *database.Database
	TableName  string
	FillStruct interface{}

	context map[string]interface{}
}

func (v *AllView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *AllView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
	return map[string]interface{}{}
}

func (v *AllView) GetContext() map[string]interface{} {
	return v.context
}

// Object sets a slice of rows from the database.
func (v *AllView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error) {
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	objects, err := v.DB.SyncQ().Select([]string{"*"}, v.TableName, dbutils.WHOutput{}, 0)
	if err != nil {
		return nil, err
	}
	fillObjects, err := v.fillObjects(objects)
	if err != nil {
		return nil, err
	}
	// CLOSE DB
	err = v.DB.Close()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{v.Name: fillObjects}, nil
}

// fillObjects fills a structure or map with data from a successful query and wraps this in a TemplateStruct.
func (v *AllView) fillObjects(objects []map[string]interface{}) ([]interface{}, error) {
	if v.FillStruct == nil {
		panic("the FillStruct field must not be nil")
	}
	var objectsStruct []interface{}
	for i := 0; i < len(objects); i++ {
		value := reflect.New(reflect.TypeOf(v.FillStruct)).Elem()
		err := dbutils.FillReflectValueFromDb(objects[i], &value)
		if err != nil {
			return nil, err
		}
		objectsStruct = append(objectsStruct, value.Interface())
	}
	return objectsStruct, nil
}

func (v *AllView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic("OnError is not implement")
}
