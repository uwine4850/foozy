package interfaces

import (
	"net/http"
)

type IMiddleware interface {
	HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager IManager, managerConfig IManagerConfig))
	RunMddl(w http.ResponseWriter, r *http.Request, manager IManager, managerConfig IManagerConfig) error
	AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager IManager, managerConfig IManagerConfig))
	RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager IManager, managerConfig IManagerConfig)
	WaitAsyncMddl()
}
