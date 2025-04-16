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
)

// ObjView displays only the HTML page only with a specific row from the database.
// Needs to be used with slug parameter URL path, specify the name of the parameter in the Slug parameter.
type ObjView struct {
	BaseView

	Name       string             `notdef:"true"`
	DB         *database.Database `notdef:"true"`
	TableName  string             `notdef:"true"`
	FillStruct interface{}        `notdef:"true"`
	Slug       string             `notdef:"true"`
}

func (v *ObjView) CloseDb() error {
	if v.DB != nil {
		if err := v.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (v *ObjView) ObjectsName() []string {
	return []string{v.Name}
}

func (v *ObjView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error) {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run ObjView object")
	if typeopr.IsPointer(v.FillStruct) {
		return nil, typeopr.ErrValueIsPointer{Value: "FillStruct"}
	}
	err := v.DB.Connect()
	if err != nil {
		return nil, err
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_DB, v.DB)

	slugValue, ok := manager.OneTimeData().GetSlugParams(v.Slug)
	if !ok {
		return nil, ErrNoSlug{v.Slug}
	}
	qRes := qb.NewSyncQB(v.DB.SyncQ()).SelectFrom("*", v.TableName).Where(
		qb.Compare(v.Slug, qb.EQUAL, slugValue),
	).Limit(1)
	qRes.Merge()
	res, err := qRes.Query()
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
		panic("the FillStruct field must not be nil")
	}
	value := reflect.New(reflect.TypeOf(v.FillStruct)).Elem()
	err := dbmapper.FillReflectValueFromDb(object, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
