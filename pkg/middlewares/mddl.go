package middlewares

import (
	"errors"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
	"sync"
)

type MddlFunc func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData)

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
func (m *Middleware) HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData)) {
	if !utils.SliceContains(m.preHandlerId, id) {
		m.preHandlerId = append(m.preHandlerId, id)
		m.preHandlerMiddlewares[id] = fn
	} else {
		m.mError = &ErrIdAlreadyExist{id}
	}
}

// AsyncHandlerMddl the middleware that will execute asynchronously before the request.
func (m *Middleware) AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData)) {
	m.asyncHandlerMiddlewares = append(m.asyncHandlerMiddlewares, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData) {
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

// SetMddlError sets an error that occurred in the middleware.
func SetMddlError(mddlErr error, manager interfaces.IManagerData) {
	manager.SetUserContext("mddlerr", mddlErr)
}

// GetMddlError get error from middleware. Used in router and in pair with SetMddlError.
func GetMddlError(manager interfaces.IManagerData) (error, error) {
	mddlErr, ok := manager.GetUserContext("mddlerr")
	if ok {
		err, ok := mddlErr.(error)
		if !ok {
			return nil, errors.New("mddlerr type is not an error")
		}
		return err, nil
	}
	return nil, nil
}

// SkipNextPage sends a command to the router to skip rendering the next page.
func SkipNextPage(manager interfaces.IManagerData) {
	manager.SetUserContext("skipNextPage", true)
}

// IsSkipNextPage checks if the page rendering should be skipped.
// The function is built into the router.
func IsSkipNextPage(manager interfaces.IManagerData) bool {
	_, ok := manager.GetUserContext("skipNextPage")
	return ok
}

// SkipNextPageAndRedirect skips the page render and redirects to another page.
func SkipNextPageAndRedirect(manager interfaces.IManagerData, w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
	SkipNextPage(manager)
}
