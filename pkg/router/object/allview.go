package object

import (
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

// AllView displays HTML page by passing all data from the selected table to it.
// If the [slug] parameter is set, all data from the table that match the condition will be output.
// If the [slug] parameter is not set, all data from the table will be output.
type AllView struct {
	BaseView
	Name       string             `notdef:"true"`
	DB         *database.Database `notdef:"true"`
	TableName  string             `notdef:"true"`
	Slug       string
	FillStruct interface{} `notdef:"true"`
}

func (v *AllView) CloseDb() error {
	if v.DB != nil {
		if err := v.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (v *AllView) ObjectsName() []string {
	return []string{v.Name}
}

// Object sets a slice of rows from the database.
func (v *AllView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run AllView object")
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_DB, v.DB)
	debug.RequestLogginIfEnable(debug.P_OBJECT, "get object from database")
	qObjects := qb.NewSyncQB(v.DB.SyncQ())
	if v.Slug != "" {
		slugValue, ok := manager.OneTimeData().GetSlugParams(v.Slug)
		if !ok {
			return nil, ErrNoSlug{v.Slug}
		}
		qObjects.SelectFrom("*", v.TableName).Where(qb.Compare(v.Slug, qb.EQUAL, slugValue))
	} else {
		qObjects.SelectFrom("*", v.TableName)
	}
	qObjects.Merge()
	objects, err := qObjects.Query()
	if err != nil {
		return nil, err
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
