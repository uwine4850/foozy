package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
)

type Handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()

type IAdapter interface {
	Adapt(pattern string, handler Handler) http.HandlerFunc
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if config.LoadedConfig().Default.Debug.Debug {
		debug.ErrorLoggingIfEnableAndWrite(w, err.Error(), err.Error())
	} else {
		debug.ErrorLoggingIfEnableAndWrite(w, err.Error(), "500 Internal server error")
	}
}

// Adapter is an object that handles router.Handler and adapts it to work as an http.HandlerFunc.
// This object triggers all the additional functionality of the handler.
// That is, it simply wraps the http.HandlerFunc controller in additional functionality and gives
// it as http.HandlerFunc.
//
// [internalErrorFunc] is responsible for handling internal errors.
// This function can be overridden using the [SetOnErrorFunc] method.
type Adapter struct {
	manager           interfaces.IManager
	middlewares       middlewares.IMiddleware
	internalErrorFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewAdapter(manager interfaces.IManager, middlewares middlewares.IMiddleware) *Adapter {
	return &Adapter{
		manager:           manager,
		middlewares:       middlewares,
		internalErrorFunc: internalServerError,
	}
}

// Adapt wraps router.Handler in additional functionality.
// It creates a new manager, starts middlewares and does other small operations.
func (a *Adapter) Adapt(pattern string, handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := debug.ClearRequestInfoLogging(); err != nil {
			a.internalErrorFunc(w, r, err)
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, fmt.Sprintf("request url: %s", r.URL))
		debug.RequestLogginIfEnable(debug.P_ROUTER, "init manager")
		newManager, err := a.newManager()
		if err != nil {
			a.internalErrorFunc(w, r, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "manager is initialized")

		// Slug params
		segments := strings.Split(strings.Trim(pattern, "/"), "/")
		urlSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		params := MatchUrlSegments(segments, urlSegments)
		if params != nil {
			newManager.OneTimeData().SetSlugParams(params)
		}

		// Run middlewares
		if skip, err := a.runPreAndAsyncMddl(w, r, newManager); err != nil {
			if a.internalErrorFunc != nil {
				a.internalErrorFunc(w, r, err)
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			} else {
				a.internalErrorFunc(w, r, err)
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			}
			return
		} else {
			if skip {
				return
			}
		}
		newManager.OneTimeData().SetUserContext(namelib.ROUTER.URL_PATTERN, pattern)
		handler(w, r, newManager)()
		a.printLog(r)
		if err := a.middlewares.RunPostMiddlewares(r, newManager); err != nil {
			if !errors.Is(err, middlewares.ErrStopMiddlewares{}) {
				a.internalErrorFunc(w, r, err)
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				return
			}
		}
	}
}

func (a *Adapter) SetOnErrorFunc(fn func(w http.ResponseWriter, r *http.Request, err error)) {
	a.internalErrorFunc = fn
}

func (a *Adapter) newManager() (interfaces.IManager, error) {
	_newManager, err := a.manager.New()
	if err != nil {
		return nil, err
	}
	newManager := _newManager.(interfaces.IManager)
	return newManager, nil
}

// runPreAndAsyncMddl runs the middleware.
// Conventional middleware is run first because it is more prioritized than asynchronous. They usually
// perform important logic, which, for example, should be run first.
// After execution of synchronous middleware, asynchronous ones are executed.
// Middleware errors and the page rendering skip algorithm are also handled here.
func (a *Adapter) runPreAndAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, error) {
	if a.middlewares != nil {
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run middlewares...")
		if err := a.middlewares.RunPreMiddlewares(w, r, manager); err != nil {
			if errors.Is(err, middlewares.ErrStopMiddlewares{}) {
				return false, nil
			}
			return false, err
		}
		if err := a.middlewares.RunAndWaitAsyncMiddlewares(w, r, manager); err != nil {
			if errors.Is(err, middlewares.ErrStopMiddlewares{}) {
				return false, nil
			}
			return false, err
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "middlewares are completed")
		// Checking the skip of the next page. Runs after a more important error check.
		if middlewares.IsSkipNextPage(manager.OneTimeData()) {
			return true, nil
		}
	}
	return false, nil
}

func (a *Adapter) printLog(request *http.Request) {
	if config.LoadedConfig().Default.Debug.PrintInfo {
		log.Printf("%s %s", request.Method, request.URL.Path)
	}
}

type RegisterHandler func(method string, pattern string, handler Handler)

type Route struct {
	Pattern  string
	Segments []string
	Handler  http.HandlerFunc
}

// Router store in itself the paths to the handler.
type Router struct {
	routes  map[string][]Route // method → slice of Route
	adapter IAdapter
}

func NewRouter(adapter IAdapter) *Router {
	return &Router{
		routes:  make(map[string][]Route),
		adapter: adapter,
	}
}

func (r *Router) HandlerSet(handlers []map[string]map[string]Handler) {
	for i := 0; i < len(handlers); i++ {
		for method, h := range handlers[i] {
			for pattern, handler := range h {
				r.Register(method, pattern, handler)
			}
		}
	}
}

func (r *Router) Register(method string, pattern string, handler Handler) {
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	adapted := r.adapter.Adapt(pattern, handler)
	r.routes[method] = append(r.routes[method], Route{
		Pattern:  pattern,
		Segments: segments,
		Handler:  adapted,
	})
}

func (r *Router) Routes() map[string][]Route {
	return r.routes
}

// ServeHTTP run handlers.
// Implementation of the [http.Handler] interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")

	if routes, ok := r.routes[method]; ok {
		for _, route := range routes {
			params := MatchUrlSegments(route.Segments, segments)
			if params != nil {
				route.Handler.ServeHTTP(w, req)
				return
			}
		}
	}

	http.NotFound(w, req)
}

func MatchUrlSegments(routeSegments, pathSegments []string) map[string]string {
	if len(routeSegments) != len(pathSegments) {
		return nil
	}

	params := make(map[string]string)

	for i := 0; i < len(routeSegments); i++ {
		rSeg := routeSegments[i]
		pSeg := pathSegments[i]

		if strings.HasPrefix(rSeg, ":") {
			paramName := rSeg[1:]
			params[paramName] = pSeg
		} else if rSeg != pSeg {
			return nil
		}
	}

	return params
}
