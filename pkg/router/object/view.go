package object

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
	"reflect"
	"strings"
)

// View The interface implements the basic structure of any View. View is used to display HTML page in a simpler and more convenient way.
// For view to work properly, you need to create a new structure (e.g. MyObjView), embed a ready view implementation
// (e.g. ObjView) into it, then initialize the selected ready ObjView, and it is important to insert the newly created
// MyObjView structure into the UserView field.
type View interface {
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func())
	Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
	OnError(e func(err error))
}

// UserView Interface for additional view control. All methods of this interface can be overridden to extend
// the functionality of the standard view. All these methods are minimally implemented in the ready view and
// do not represent useful functionality, to add functionality they need to be overridden.
// The Context method adds new variables to the HTML template.
// The Permissions method performs the returned function and returns false.
type UserView interface {
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func())
}

// TemplateStruct The structure is used to display the table data in the view.
// The data can be accessed using the F method, which takes a field name.
// It is desirable to always write the field name with a small letter (even public fields of a structure) to avoid
// separating data in the form of a structure and data in the form of a map.
type TemplateStruct struct {
	s reflect.Value
	m map[string]string
}

// F outputs the field data by its name.
func (t TemplateStruct) F(name string) string {
	if t.s.Kind() != reflect.Invalid {
		return t.s.FieldByName(strings.ToUpper(string(name[0])) + name[1:]).String()
	}
	if t.m != nil {
		return t.m[name]
	}
	return ""
}
