package router

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
)

// RedirectError redirect to the page and provide error information for it.
func RedirectError(w http.ResponseWriter, r *http.Request, path string, err string, manager interfaces.IManager) {
	manager.SetUserContext("error", err)
	http.Redirect(w, r, path, http.StatusFound)
}

// HandleRedirectError handling by the template engine of an error sent by the HandleRedirectError function.
// In the template you can get an error using the error variable.
func HandleRedirectError(manager interfaces.IManagerData) {
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": ""})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string)})
		manager.DelUserContext("error")
	}
}

// ServerError displaying a 500 error to the user.
func ServerError(w http.ResponseWriter, error string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(error))
}

// ServerForbidden displaying a 403 error to the user.
func ServerForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 forbidden"))
}
