package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type MultipleObject struct {
	Name       string
	TaleName   string
	SlugName   string
	AIField    string
	FillStruct interface{}
}

type MultipleObjectView struct {
	IView

	DB              *database.Database
	MultipleObjects []MultipleObject

	context map[string]interface{}
}

func (v *MultipleObjectView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *MultipleObjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
	return map[string]interface{}{}
}

func (v *MultipleObjectView) GetContext() map[string]interface{} {
	return v.context
}

func (v *MultipleObjectView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error) {
	v.context = make(map[string]interface{})
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(v.MultipleObjects); i++ {
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.MultipleObjects[i].SlugName)
		if !ok {
			return nil, ErrNoSlug{v.MultipleObjects[i].SlugName}
		}
		res, err := v.DB.SyncQ().Select([]string{"*"}, v.MultipleObjects[i].TaleName, dbutils.WHEquals(map[string]interface{}{
			v.MultipleObjects[i].AIField: slugValue,
		}, "AND"), 1)
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, ErrNoData{}
		}
		value, err := v.fillObject(res[0], v.MultipleObjects[i].FillStruct)
		if err != nil {
			return nil, err
		}
		v.context[v.MultipleObjects[i].Name] = value.Interface()
	}
	// CLOSE DB
	err = v.DB.Close()
	if err != nil {
		return nil, err
	}
	return v.context, nil
}

func (v *MultipleObjectView) fillObject(object map[string]interface{}, fillStruct interface{}) (*reflect.Value, error) {
	if fillStruct == nil {
		panic("the FillStruct field must not be nil")
	}
	value := reflect.New(reflect.TypeOf(fillStruct)).Elem()
	err := dbutils.FillReflectValueFromDb(object, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (v *MultipleObjectView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic("OnError is not implement")
}
