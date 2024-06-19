package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

// RedirectError redirect to the page and provide error information for it.
func RedirectError(w http.ResponseWriter, r *http.Request, path string, _err string, manager interfaces.IManager) {
	uval := url.Values{}
	uval.Add(namelib.REDIRECT_ERROR, _err)
	newUrl := fmt.Sprintf("%s?%s", path, uval.Encode())
	http.Redirect(w, r, newUrl, http.StatusFound)
	debug.ErrorLogginIfEnable(_err, manager.Config())
}

// CatchRedirectError handling by the template engine of an error sent by the CatchRedirectError function.
// In the template you can get an error using the error variable.
func CatchRedirectError(r *http.Request, manager interfaces.IManager) {
	q := r.URL.Query()
	redirectError := q.Get(namelib.REDIRECT_ERROR)
	if redirectError != "" {
		manager.Render().SetContext(map[string]interface{}{namelib.REDIRECT_ERROR: redirectError})
	}
}

// ServerError displaying a 500 error to the user.
func ServerError(w http.ResponseWriter, error string, manager interfaces.IManagerConfig) {
	w.WriteHeader(http.StatusInternalServerError)
	if manager.IsDebug() {
		debug.ErrorLoggingIfEnableAndWrite(w, []byte(error), manager)
	} else {
		debug.ErrorLoggingIfEnableAndWrite(w, []byte("500 Internal server error"), manager)
	}
}

// ServerForbidden displaying a 403 error to the user.
func ServerForbidden(w http.ResponseWriter, manager interfaces.IManagerConfig) {
	w.WriteHeader(http.StatusForbidden)
	debug.ErrorLoggingIfEnableAndWrite(w, []byte("403 forbidden"), manager)
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
