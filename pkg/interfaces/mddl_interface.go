package interfaces

import (
	"net/http"
)

type IMiddleware interface {
	HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	RunMddl(w http.ResponseWriter, r *http.Request, manager IManager) error
	AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager IManager)
	WaitAsyncMddl()
}
