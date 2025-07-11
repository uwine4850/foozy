package object

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

type Context map[string]interface{}

// IView the interface implements the basic structure of any IView. ITemplateView is used to display HTML page in a simpler and more convenient way.
// For the view to work correctly, you need to create a new structure (for example MyObjView), embed a ready-made implementation of the view
// (for example ObjView) into it, then you need to initialize this structure in the ITemplateView field in the TemplateView data type.
type IView interface {
	// Object receives data from the selected table and writes it to a variable structure.
	// IMPORTANT: connects to the database in this method (or others), but closes the connection only in the TemplateView.
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error)
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (bool, func())
	OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error)
	ObjectsName() []string
}

// IViewDatabase interface that provides object unified access to the database.
// Since each object queries the database, it is necessary to unify access
// to the database so as not to be dependent on a particular database.
type IViewDatabase interface {
	SelectAll(tableName string) ([]map[string]interface{}, error)
	SelectWhereEqual(tableName string, colName string, val any) ([]map[string]interface{}, error)
}

// ViewMysqlDatabase implementation of the [IViewDatabase] interface for a MySql database.
type ViewMysqlDatabase struct {
	db interfaces.DatabaseInteraction
}

func NewViewMysqlDatabase(db interfaces.DatabaseInteraction) *ViewMysqlDatabase {
	return &ViewMysqlDatabase{
		db: db,
	}
}

func (d *ViewMysqlDatabase) SelectAll(tableName string) ([]map[string]interface{}, error) {
	return d.db.SyncQ().Query(fmt.Sprintf("SELECT * FROM %s", tableName))
}

func (d *ViewMysqlDatabase) SelectWhereEqual(tableName string, colName string, val any) ([]map[string]interface{}, error) {
	return d.db.SyncQ().Query(fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, colName), val)
}

type BaseView struct{}

func (v *BaseView) CloseDb() error {
	panic("CloseDb is not implement. Please implement this method in your structure.")
}

func (v *BaseView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error) {
	return Context{}, nil
}

func (v *BaseView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error) {
	return Context{}, nil
}

func (v *BaseView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (bool, func()) {
	return true, func() {}
}

func (v *BaseView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
	panic("OnError is not implement. Please implement this method in your structure.")
}

// GetContext retrieves the Context from the manager.
// It is important to understand that this method can only be used when the IView.Object method has completed running,
// for example in IView.Context.
func GetContext(manager interfaces.Manager) (Context, error) {
	objectInterface, ok := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	if !ok {
		return nil, errors.New("unable to get object context")
	}
	object := objectInterface.(Context)
	return object, nil
}

type ErrNoSlug struct {
	SlugName string
}

func (e ErrNoSlug) Error() string {
	return fmt.Sprintf("slug parameter %s not found", e.SlugName)
}

type ErrNoData struct {
}

func (e ErrNoData) Error() string {
	return "no data to display was found"
}
