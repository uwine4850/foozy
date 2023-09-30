package middlewares

import (
	"github.com/uwine4850/foozy/internal/interfaces"
	"github.com/uwine4850/foozy/internal/utils"
	"net/http"
	"sync"
)

type MddlFunc func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)

type Middleware struct {
	preHandlerMiddlewares   map[int]MddlFunc
	preHandlerId            []int
	asyncHandlerMiddlewares []MddlFunc
	Context                 sync.Map
	mError                  error
	wg                      sync.WaitGroup
}

func NewMiddleware() *Middleware {
	return &Middleware{preHandlerMiddlewares: make(map[int]MddlFunc)}
}

// HandlerMddl the middleware that will be executed before the request handler.
// id indicates the order of execution of the current middleware. No two identical id's can be created.
func (m *Middleware) HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	if !utils.SliceContains(m.preHandlerId, id) {
		m.preHandlerId = append(m.preHandlerId, id)
		m.preHandlerMiddlewares[id] = fn
	} else {
		m.mError = &ErrIdAlreadyExist{id}
	}
}

// AsyncHandlerMddl the middleware that will execute asynchronously before the request.
func (m *Middleware) AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	m.asyncHandlerMiddlewares = append(m.asyncHandlerMiddlewares, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		defer m.wg.Done()
		fn(w, r, manager)
	})
}

// RunMddl running synchronous middleware.
func (m *Middleware) RunMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	if m.mError != nil {
		return m.mError
	}
	for _, handlerFunc := range m.preHandlerMiddlewares {
		handlerFunc(w, r, manager)
	}
	return nil
}

// RunAsyncMddl running asynchronous middleware.
// IMPORTANT: You must run the WaitAsyncMddl method at the selected location to complete correctly.
func (m *Middleware) RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	for i := 0; i < len(m.asyncHandlerMiddlewares); i++ {
		m.wg.Add(1)
		go m.asyncHandlerMiddlewares[i](w, r, manager)
	}
}

// WaitAsyncMddl waits for the execution of all asynchronous middleware.
func (m *Middleware) WaitAsyncMddl() {
	m.wg.Wait()
}
