package middlewares

import (
	"errors"
	"net/http"
	"sync"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/utils"
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

// SetMddlError sets an error that occurred in the middleware.
func SetMddlError(mddlErr error, manager interfaces.IManagerOneTimeData) {
	manager.SetUserContext(namelib.MDDL_ERROR, mddlErr)
}

// GetMddlError get error from middleware. Used in router and in pair with SetMddlError.
func GetMddlError(manager interfaces.IManagerOneTimeData) (error, error) {
	mddlErr, ok := manager.GetUserContext(namelib.MDDL_ERROR)
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
func SkipNextPage(manager interfaces.IManagerOneTimeData) {
	manager.SetUserContext(namelib.SKIP_NEXT_PAGE, true)
}

// IsSkipNextPage checks if the page rendering should be skipped.
// The function is built into the router.
func IsSkipNextPage(manager interfaces.IManagerOneTimeData) bool {
	_, ok := manager.GetUserContext(namelib.SKIP_NEXT_PAGE)
	return ok
}

// SkipNextPageAndRedirect skips the page render and redirects to another page.
func SkipNextPageAndRedirect(manager interfaces.IManagerOneTimeData, w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
	SkipNextPage(manager)
}
