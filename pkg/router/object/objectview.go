package object

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

// ObjView displays only the HTML page only with a specific row from the database.
// Needs to be used with slug parameter URL path, specify the name of the parameter in the Slug parameter.
type ObjView struct {
	BaseView

	Name       string        `notdef:"true"`
	TableName  string        `notdef:"true"`
	Database   IViewDatabase `notdef:"true"`
	FillStruct interface{}   `notdef:"true"`
	Slug       string        `notdef:"true"`
}

func (v *ObjView) ObjectsName() []string {
	return []string{v.Name}
}

func (v *ObjView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run ObjView object")
	if typeopr.IsPointer(v.FillStruct) {
		return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
	}

	slugValue, ok := manager.OneTimeData().GetSlugParams(v.Slug)
	if !ok {
		return nil, ErrNoSlug{v.Slug}
	}
	res, err := v.Database.SelectWhereEqual(v.TableName, v.Slug, slugValue)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrNoData{}
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "fill object")
	value, err := v.fillObject(res[0])
	if err != nil {
		return nil, err
	}
	return Context{v.Name: value.Interface()}, nil
}

func (v *ObjView) fillObject(object map[string]interface{}) (*reflect.Value, error) {
	if v.FillStruct == nil {
		return nil, errors.New("the FillStruct field must not be nil")
	}
	fillType := reflect.TypeOf(v.FillStruct)
	value := reflect.New(fillType).Elem()
	err := mapper.FillStructFromDb(&value, &object)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
