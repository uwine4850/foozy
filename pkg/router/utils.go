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
	uval.Add(namelib.ROUTER.REDIRECT_ERROR, _err)
	newUrl := fmt.Sprintf("%s?%s", path, uval.Encode())
	http.Redirect(w, r, newUrl, http.StatusFound)
	debug.ErrorLogginIfEnable(_err, manager.Config())
}

// CatchRedirectError handling by the template engine of an error sent by the CatchRedirectError function.
// In the template you can get an error using the error variable.
func CatchRedirectError(r *http.Request, manager interfaces.IManager) {
	q := r.URL.Query()
	redirectError := q.Get(namelib.ROUTER.REDIRECT_ERROR)
	if redirectError != "" {
		manager.Render().SetContext(map[string]interface{}{namelib.ROUTER.REDIRECT_ERROR: redirectError})
	}
}

// ServerError displaying a 500 error to the user.
func ServerError(w http.ResponseWriter, error string, manager interfaces.IManager) {
	manager.OneTimeData().SetUserContext(namelib.ROUTER.SERVER_ERROR, error)
	w.WriteHeader(http.StatusInternalServerError)
	if manager.Config().IsDebug() {
		debug.ErrorLoggingIfEnableAndWrite(w, []byte(error), manager.Config())
	} else {
		debug.ErrorLoggingIfEnableAndWrite(w, []byte("500 Internal server error"), manager.Config())
	}
}

// ServerForbidden displaying a 403 error to the user.
func ServerForbidden(w http.ResponseWriter, manager interfaces.IManager) {
	manager.OneTimeData().SetUserContext(namelib.ROUTER.SERVER_FORBIDDEN_ERROR, "403 forbidden")
	w.WriteHeader(http.StatusForbidden)
	debug.ErrorLoggingIfEnableAndWrite(w, []byte("403 forbidden"), manager.Config())
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
