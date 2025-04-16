package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

type MddlFunc func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
type AsyncMddlFunc func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, wg *sync.WaitGroup)

type Middleware struct {
	syncMiddlewares  map[int]MddlFunc
	syncHandlerId    []int
	asyncMiddlewares []AsyncMddlFunc
	Context          sync.Map
	mError           error
}

func NewMiddleware() *Middleware {
	return &Middleware{syncMiddlewares: make(map[int]MddlFunc)}
}

// SyncMddl the middleware that will be executed before the request handler.
// id indicates the order of execution of the current middleware. No two identical id's can be created.
func (m *Middleware) SyncMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	if !fslice.SliceContains(m.syncHandlerId, id) {
		m.syncHandlerId = append(m.syncHandlerId, id)
		m.syncMiddlewares[id] = fn
	} else {
		m.mError = &ErrIdAlreadyExist{id}
	}
}

// AsyncHandlerMddl the middleware that will execute asynchronously before the request.
func (m *Middleware) AsyncMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	m.asyncMiddlewares = append(m.asyncMiddlewares, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, wg *sync.WaitGroup) {
		defer wg.Done()
		fn(w, r, manager)
	})
}

// RunMddl running synchronous middleware.
func (m *Middleware) RunMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	if m.mError != nil {
		return m.mError
	}
	sort.Ints(m.syncHandlerId)
	for i := 0; i < len(m.syncHandlerId); i++ {
		debug.RequestLogginIfEnable(debug.P_MIDDLEWARE, fmt.Sprintf("sync middleware with id %s is running", strconv.Itoa(i)))
		handlerFunc := m.syncMiddlewares[i]
		handlerFunc(w, r, manager)
	}
	return nil
}

// RunAsyncMddl running asynchronous middleware.
// IMPORTANT: You must run the WaitAsyncMddl method at the selected location to complete correctly.
func (m *Middleware) RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, wg *sync.WaitGroup) {
	for i := 0; i < len(m.asyncMiddlewares); i++ {
		debug.RequestLogginIfEnable(debug.P_MIDDLEWARE, fmt.Sprintf("async middleware with id %s is running", strconv.Itoa(i)))
		wg.Add(1)
		go m.asyncMiddlewares[i](w, r, manager, wg)
	}
}

// SetMddlError sets an error that occurred in the middleware.
func SetMddlError(mddlErr error, manager interfaces.IManagerOneTimeData) {
	manager.SetUserContext(namelib.ROUTER.MDDL_ERROR, mddlErr)
	if config.LoadedConfig().Default.Debug.RequestInfoLog && config.LoadedConfig().Default.Debug.RequestInfoLogPath == "" {
		panic("unable to create request info log file. File path not set")
	}
	debug.WriteLog(config.LoadedConfig().Default.Debug.SkipLoggingLevel-1, config.LoadedConfig().Default.Debug.RequestInfoLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, debug.P_MIDDLEWARE, fmt.Sprintf("middleware error: %s", mddlErr.Error()))
}

// GetMddlError get error from middleware. Used in router and in pair with SetMddlError.
func GetMddlError(manager interfaces.IManagerOneTimeData) (error, error) {
	mddlErr, ok := manager.GetUserContext(namelib.ROUTER.MDDL_ERROR)
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
	manager.SetUserContext(namelib.ROUTER.SKIP_NEXT_PAGE, true)
	urlPattern, _ := manager.GetUserContext(namelib.ROUTER.URL_PATTERN)
	debug.RequestLogginIfEnable(debug.P_MIDDLEWARE, fmt.Sprintf("skip page at %s", urlPattern))
}

// IsSkipNextPage checks if the page rendering should be skipped.
// The function is built into the router.
func IsSkipNextPage(manager interfaces.IManagerOneTimeData) bool {
	_, ok := manager.GetUserContext(namelib.ROUTER.SKIP_NEXT_PAGE)
	return ok
}

// SkipNextPageAndRedirect skips the page render and redirects to another page.
func SkipNextPageAndRedirect(manager interfaces.IManagerOneTimeData, w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
	SkipNextPage(manager)
	debug.RequestLogginIfEnable(debug.P_MIDDLEWARE, fmt.Sprintf("redirect to %s", path))
}
