package object

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
)

type ObjectContext map[string]interface{}

// IView The interface implements the basic structure of any IView. ITemplateView is used to display HTML page in a simpler and more convenient way.
// For the view to work correctly, you need to create a new structure (for example MyObjView), embed a ready-made implementation of the view
// (for example ObjView) into it, then you need to initialize this structure in the ITemplateView field in the TemplateView data type.
type IView interface {
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) ObjectContext
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func())
	OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)
}
