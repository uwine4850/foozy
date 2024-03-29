package router

import (
	"encoding/json"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
)

// RedirectError redirect to the page and provide error information for it.
func RedirectError(w http.ResponseWriter, r *http.Request, path string, err string, manager interfaces.IManager) {
	manager.SetUserContext("error", err)
	http.Redirect(w, r, path, http.StatusFound)
}

// CatchRedirectError handling by the template engine of an error sent by the CatchRedirectError function.
// In the template you can get an error using the error variable.
func CatchRedirectError(manager interfaces.IManagerData) {
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

// SendJson sends json-formatted data to the page.
func SendJson(data interface{}, w http.ResponseWriter) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}
