package router

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
)

func RedirectError(w http.ResponseWriter, r *http.Request, path string, err string, manager interfaces.IManager) {
	manager.SetUserContext("error", err)
	http.Redirect(w, r, path, http.StatusFound)
}

func ServerError(w http.ResponseWriter, error string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(error))
}

func ServerForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 forbidden"))
}
