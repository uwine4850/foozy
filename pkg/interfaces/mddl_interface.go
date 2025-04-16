package interfaces

import (
	"net/http"
	"sync"
)

type IMiddleware interface {
	SyncMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	AsyncMddl(fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	RunMddl(w http.ResponseWriter, r *http.Request, manager IManager) error
	RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager IManager, wg *sync.WaitGroup)
}
