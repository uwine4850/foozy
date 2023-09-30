package interfaces

import (
	"net/http"
)

type IMiddleware interface {
	PreHandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	RunPreMddl(w http.ResponseWriter, r *http.Request, manager IManager) error
}
