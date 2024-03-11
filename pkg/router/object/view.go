package object

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
)

// View The interface implements the basic structure of any View. View is used to display HTML page in a simpler and more convenient way.
// For the view to work correctly, you need to create a new structure (for example MyObjView), embed a ready-made implementation of the view
// (for example ObjView) into it, then you need to initialize this structure in the View field in the TemplateView data type.
type View interface {
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func())
	Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
	OnError(e func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error))
}
