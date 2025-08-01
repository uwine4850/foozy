package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fstruct"
)

type MultipleObject struct {
	Name       string      `notdef:"true"`
	TableName  string      `notdef:"true"`
	SlugName   string      `notdef:"true"`
	SlugField  string      `notdef:"true"`
	FillStruct interface{} `notdef:"true"`
}

type MultipleObjectView struct {
	BaseView

	Database        IViewDatabase    `notdef:"true"`
	MultipleObjects []MultipleObject `notdef:"true"`
}

func (v *MultipleObjectView) ObjectsName() []string {
	names := []string{}
	for i := 0; i < len(v.MultipleObjects); i++ {
		names = append(names, v.MultipleObjects[i].Name)
	}
	return names
}

func (v *MultipleObjectView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run MultipleObjectView object")
	if err := v.checkMultipleObject(); err != nil {
		return nil, err
	}
	context := make(Context)
	debug.RequestLogginIfEnable(debug.P_OBJECT, "start fill objects")
	for i := 0; i < len(v.MultipleObjects); i++ {
		if typeopr.IsPointer(v.MultipleObjects[i].FillStruct) {
			return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
		}
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.MultipleObjects[i].SlugName)
		if !ok {
			return nil, ErrNoSlug{v.MultipleObjects[i].SlugName}
		}
		res, err := v.Database.SelectWhereEqual(v.MultipleObjects[i].TableName, v.MultipleObjects[i].SlugField, slugValue)
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
