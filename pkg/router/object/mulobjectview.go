package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbmapper"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type MultipleObject struct {
	Name       string
	TaleName   string
	SlugName   string
	SlugField  string
	FillStruct interface{}
}

type MultipleObjectView struct {
	BaseView

	DB              *database.Database
	MultipleObjects []MultipleObject
}

func (v *MultipleObjectView) CloseDb() error {
	if v.DB != nil {
		if err := v.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (v *MultipleObjectView) ObjectsName() []string {
	names := []string{}
	for i := 0; i < len(v.MultipleObjects); i++ {
		names = append(names, v.MultipleObjects[i].Name)
	}
	return names
}

func (v *MultipleObjectView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
	context := make(ObjectContext)
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_DB, v.DB)

	for i := 0; i < len(v.MultipleObjects); i++ {
		if typeopr.IsPointer(v.MultipleObjects[i].FillStruct) {
			return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
		}
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
	err := dbmapper.FillReflectValueFromDb(object, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
