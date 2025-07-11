package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

// AllView displays HTML page by passing all data from the selected table to it.
// If the [slug] parameter is set, all data from the table that match the condition will be output.
// If the [slug] parameter is not set, all data from the table will be output.
type AllView struct {
	BaseView
	Name       string        `notdef:"true"`
	TableName  string        `notdef:"true"`
	Database   IViewDatabase `notdef:"true"`
	Slug       string
	FillStruct interface{} `notdef:"true"`
}

func (v *AllView) ObjectsName() []string {
	return []string{v.Name}
}

// Object sets a slice of rows from the database.
func (v *AllView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run AllView object")
	debug.RequestLogginIfEnable(debug.P_OBJECT, "get object from database")
	var objects []map[string]interface{}
	if v.Slug != "" {
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.Slug)
		if !ok {
			return nil, ErrNoSlug{v.Slug}
		}
		res, err := v.Database.SelectWhereEqual(v.TableName, v.Slug, slugValue)
		if err != nil {
			return nil, err
		}
		objects = res
	} else {
		res, err := v.Database.SelectAll(v.TableName)
		if err != nil {
			return nil, err
		}
		objects = res
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "fill objects")
	fillObjects, err := v.fillObjects(objects)
	if err != nil {
		return nil, err
	}
	return Context{v.Name: fillObjects}, nil
}

// fillObjects fills a structure or map with data from a successful query and wraps this in a TemplateStruct.
func (v *AllView) fillObjects(objects []map[string]interface{}) ([]interface{}, error) {
	if v.FillStruct == nil {
		panic("the FillStruct field must not be nil")
	}
	if typeopr.IsPointer(v.FillStruct) {
		return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
	}
	var objectsStruct []interface{}
	for i := 0; i < len(objects); i++ {
		fillType := reflect.TypeOf(v.FillStruct)
		value := reflect.New(fillType).Elem()
		err := mapper.FillStructFromDb(&value, &objects[i])
		if err != nil {
			return nil, err
		}
		objectsStruct = append(objectsStruct, value.Interface())
	}
	return objectsStruct, nil
}
