package object

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

type Context map[string]interface{}

// IView The interface implements the basic structure of any IView. ITemplateView is used to display HTML page in a simpler and more convenient way.
// For the view to work correctly, you need to create a new structure (for example MyObjView), embed a ready-made implementation of the view
// (for example ObjView) into it, then you need to initialize this structure in the ITemplateView field in the TemplateView data type.
type IView interface {
	// Object receives data from the selected table and writes it to a variable structure.
	// IMPORTANT: connects to the database in this method (or others), but closes the connection only in the TemplateView.
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error)
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func())
	OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)
	ObjectsName() []string
}

type BaseView struct{}

func (v *BaseView) CloseDb() error {
	panic("CloseDb is not implement. Please implement this method in your structure.")
}

func (v *BaseView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error) {
	return Context{}, nil
}

func (v *BaseView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (Context, error) {
	return Context{}, nil
}

func (v *BaseView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *BaseView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic("OnError is not implement. Please implement this method in your structure.")
}

// GetContext retrieves the Context from the manager.
// It is important to understand that this method can only be used when the IView.Object method has completed running,
// for example in IView.Context.
func GetContext(manager interfaces.IManager) (Context, error) {
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
