package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/config"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
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

	MultipleObjects []MultipleObject `notdef:"true"`
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
	dbRead, err := manager.Database().ConnectionPool(config.LoadedConfig().Default.Database.MainConnectionPoolName)
	if err != nil {
		return nil, err
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "start fill objects")
	for i := 0; i < len(v.MultipleObjects); i++ {
		if typeopr.IsPointer(v.MultipleObjects[i].FillStruct) {
			return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
		}
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.MultipleObjects[i].SlugName)
		if !ok {
			return nil, ErrNoSlug{v.MultipleObjects[i].SlugName}
		}
		qRes := qb.NewSyncQB(dbRead.SyncQ()).SelectFrom("*", v.MultipleObjects[i].TaleName).Where(
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
	fillType := reflect.TypeOf(fillStruct)
	value := reflect.New(fillType).Elem()
	err := mapper.FillStructFromDb(&value, &object)
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
