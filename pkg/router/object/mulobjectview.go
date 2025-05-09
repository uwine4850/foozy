package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbmapper"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fstruct"
)

type MultipleObject struct {
	Name       string      `notdef:"true"`
	TaleName   string      `notdef:"true"`
	SlugName   string      `notdef:"true"`
	SlugField  string      `notdef:"true"`
	FillStruct interface{} `notdef:"true"`
}

type MultipleObjectView struct {
	BaseView

	DB              *database.Database `notdef:"true"`
	MultipleObjects []MultipleObject   `notdef:"true"`
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

func (v *MultipleObjectView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run MultipleObjectView object")
	if err := v.checkMultipleObject(); err != nil {
		return nil, err
	}
	context := make(Context)
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_DB, v.DB)
	debug.RequestLogginIfEnable(debug.P_OBJECT, "start fill objects")
	for i := 0; i < len(v.MultipleObjects); i++ {
		if typeopr.IsPointer(v.MultipleObjects[i].FillStruct) {
			return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
		}
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.MultipleObjects[i].SlugName)
		if !ok {
			return nil, ErrNoSlug{v.MultipleObjects[i].SlugName}
		}
		qRes := qb.NewSyncQB(v.DB.SyncQ()).SelectFrom("*", v.MultipleObjects[i].TaleName).Where(
			qb.Compare(v.MultipleObjects[i].SlugField, qb.EQUAL, slugValue),
		).Limit(1)
		qRes.Merge()
		res, err := qRes.Query()
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

func (v *MultipleObjectView) checkMultipleObject() error {
	for i := 0; i < len(v.MultipleObjects); i++ {
		if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&v.MultipleObjects[i])); err != nil {
			return err
		}
	}
	return nil
}
