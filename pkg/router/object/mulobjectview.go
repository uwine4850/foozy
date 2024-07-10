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
	SlugField  string
	FillStruct interface{}
}

type MultipleObjectView struct {
	IView

	DB              *database.Database
	MultipleObjects []MultipleObject
}

func (v *MultipleObjectView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *MultipleObjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) ObjectContext {
	return ObjectContext{}
}

func (v *MultipleObjectView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
	context := make(ObjectContext)
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = v.DB.Close()
		if err != nil {
			v.OnError(w, r, manager, err)
		}
	}()

	for i := 0; i < len(v.MultipleObjects); i++ {
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.MultipleObjects[i].SlugName)
		if !ok {
			return nil, ErrNoSlug{v.MultipleObjects[i].SlugName}
		}
		res, err := v.DB.SyncQ().Select([]string{"*"}, v.MultipleObjects[i].TaleName, dbutils.WHEquals(map[string]interface{}{
			v.MultipleObjects[i].SlugField: slugValue,
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
		context[v.MultipleObjects[i].Name] = value.Interface()
	}
	return context, nil
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
	panic("OnError is not implement. Please implement this method in your structure.")
}
