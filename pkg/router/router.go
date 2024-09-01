package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
)

type Handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()

// muxRouter represents a single URL handler that can fire method handlers according to those sent by the client.
type muxRouter struct {
	Get     Handler
	Post    Handler
	Put     Handler
	Delete  Handler
	Options Handler
	Ws      Handler
}

type Router struct {
	mux        http.ServeMux
	routes     map[string]muxRouter
	manager    interfaces.IManager
	middleware interfaces.IMiddleware
}

func NewRouter(manager interfaces.IManager) *Router {
	return &Router{mux: *http.NewServeMux(), manager: manager, routes: map[string]muxRouter{}}
}

func (rt *Router) GetMux() *http.ServeMux {
	return &rt.mux
}

// RegisterAll registers all route handlers
func (rt *Router) RegisterAll() {
	rt.registerAllHandlers()
}

func (rt *Router) Get(pattern string, handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()) {
	_muxRouter := rt.getMuxRouter(pattern)
	if _muxRouter.Get != nil {
		panic(fmt.Sprintf("the %s method on the %s path is already mounted", "GET", pattern))
	}
	_muxRouter.Get = handler
	rt.setMuxRouter(pattern, _muxRouter)
}

func (rt *Router) Post(pattern string, handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()) {
	_muxRouter := rt.getMuxRouter(pattern)
	if _muxRouter.Post != nil {
		panic(fmt.Sprintf("the %s method on the %s path is already mounted", "POST", pattern))
	}
	_muxRouter.Post = handler
	rt.setMuxRouter(pattern, _muxRouter)
}

func (rt *Router) Put(pattern string, handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()) {
	_muxRouter := rt.getMuxRouter(pattern)
	if _muxRouter.Put != nil {
		panic(fmt.Sprintf("the %s method on the %s path is already mounted", "PUT", pattern))
	}
	_muxRouter.Put = handler
	rt.setMuxRouter(pattern, _muxRouter)
}

func (rt *Router) Delete(pattern string, handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()) {
	_muxRouter := rt.getMuxRouter(pattern)
	if _muxRouter.Delete != nil {
		panic(fmt.Sprintf("the %s method on the %s path is already mounted", "DELETE", pattern))
	}
	_muxRouter.Delete = handler
	rt.setMuxRouter(pattern, _muxRouter)
}

func (rt *Router) Options(pattern string, handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()) {
	_muxRouter := rt.getMuxRouter(pattern)
	if _muxRouter.Options != nil {
		panic(fmt.Sprintf("the %s method on the %s path is already mounted", "OPTIONS", pattern))
	}
	_muxRouter.Options = handler
	rt.setMuxRouter(pattern, _muxRouter)
}

func (rt *Router) Ws(pattern string, handler func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()) {
	_muxRouter := rt.getMuxRouter(pattern)
	if _muxRouter.Ws != nil {
		panic(fmt.Sprintf("the %s method on the %s path is already mounted", "WS", pattern))
	}
	_muxRouter.Ws = handler
	rt.setMuxRouter(pattern, _muxRouter)
}

func (rt *Router) SetMiddleware(middleware interfaces.IMiddleware) {
	rt.middleware = middleware
}

// getMuxRouter returns a muxRouter structure.
// If it does not exist, creates and returns it.
func (rt *Router) getMuxRouter(pattern string) muxRouter {
	_muxRouter, exists := rt.routes[pattern]
	if !exists {
		_muxRouter = muxRouter{}
	}
	return _muxRouter
}

// setMuxRouter set the muxRouter structure in the routes map.
// Works well in conjunction with the getMuxRouter method.
func (rt *Router) setMuxRouter(pattern string, _muxRouter muxRouter) {
	rt.routes[pattern] = _muxRouter
}

// validateMethod checks whether the method is allowed to be applied on the given URL.
func (rt *Router) validateMethod(handler Handler, method string, w http.ResponseWriter) bool {
	if handler == nil {
		w.Header().Set("Allow", method)
		http.Error(w, fmt.Sprintf("Method %s Not Allowed", method), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// register passed to muxRouter.
// When navigating to a URL, the handler will look for a function from muxRouter to run that handler.
// Various services are also started here for the correct operation of the processor.
func (rt *Router) register(_muxRouter muxRouter, urlPattern string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := rt.initManager(); err != nil {
			ServerError(writer, err.Error(), rt.manager)
			return
		}
		rt.manager.OneTimeData().SetUserContext(namelib.ROUTER.URL_PATTERN, urlPattern)

		parseUrl := ParseSlugIndex(fstring.SplitUrl(urlPattern))
		if request.URL.Path != "/" && len(parseUrl) > 0 {
			res, params := HandleSlugUrls(parseUrl, fstring.SplitUrl(urlPattern), fstring.SplitUrl(request.URL.Path))
			if res != request.URL.Path {
				http.NotFound(writer, request)
				return
			}
			if params != nil {
				rt.manager.OneTimeData().SetSlugParams(params)
			}
		}

		// Run middlewares.
		if skip, err := rt.runMddl(writer, request); err != nil {
			ServerError(writer, err.Error(), rt.manager)
			return
		} else {
			if skip {
				return
			}
		}

		connection := request.Header.Get("Connection")
		if connection != "" && connection == "Upgrade" {
			handler := _muxRouter.Ws
			if !rt.validateMethod(handler, "WS", writer) {
				return
			}
			handler(writer, request, rt.manager)()
			return
		}
		switch request.Method {
		case http.MethodGet:
			handler := _muxRouter.Get
			if !rt.validateMethod(handler, "GET", writer) {
				return
			}
			handler(writer, request, rt.manager)()
		case http.MethodPost:
			handler := _muxRouter.Post
			if !rt.validateMethod(handler, "POST", writer) {
				return
			}
			handler(writer, request, rt.manager)()
		case http.MethodPut:
			handler := _muxRouter.Put
			if !rt.validateMethod(handler, "PUT", writer) {
				return
			}
			handler(writer, request, rt.manager)()
		case http.MethodDelete:
			handler := _muxRouter.Delete
			if !rt.validateMethod(handler, "DELETE", writer) {
				return
			}
			handler(writer, request, rt.manager)()
		case http.MethodOptions:
			handler := _muxRouter.Options
			if !rt.validateMethod(handler, "OPTIONS", writer) {
				return
			}
			handler(writer, request, rt.manager)()
		}
		rt.printLog(request)
	}
}

// initManager initializes a new manager instance.
// Must be called for each new request.
func (rt *Router) initManager() error {
	if err := manager.CreateAndSetNewManagerData(rt.manager); err != nil {
		return err
	}
	if rt.manager.Render() != nil {
		if err := tmlengine.CreateAndSetNewRenderInstance(rt.manager); err != nil {
			return err
		}
	}
	return nil
}

// runMddl runs the middleware.
// Conventional middleware is run first because it is more prioritized than asynchronous. They usually
// perform important logic, which, for example, should be run first.
// After execution of synchronous middleware, asynchronous ones are executed.
// Middleware errors and the page rendering skip algorithm are also handled here.
func (rt *Router) runMddl(w http.ResponseWriter, r *http.Request) (bool, error) {
	if rt.middleware != nil {
		// Running synchronous middleware.
		err := rt.middleware.RunMddl(w, r, rt.manager)
		if err != nil {
			return false, err
		}
		// Running asynchronous middleware.
		rt.middleware.RunAsyncMddl(w, r, rt.manager)
		// Waiting for all asynchronous middleware to complete.
		rt.middleware.WaitAsyncMddl()
		// Handling middleware errors.
		mddlErr, err := middlewares.GetMddlError(rt.manager.OneTimeData())
		if err != nil {
			return false, err
		}
		if mddlErr != nil {
			return false, errors.New(mddlErr.Error())
		}
		// Checking the skip of the next page. Runs after a more important error check.
		if middlewares.IsSkipNextPage(rt.manager.OneTimeData()) {
			return true, nil
		}
	}
	return false, nil
}

func (rt *Router) registerAllHandlers() {
	for pattern, _muxRouter := range rt.routes {
		rt.mux.Handle(splitUrlFromFirstSlug(pattern), rt.register(_muxRouter, pattern))
	}
}

func (rt *Router) printLog(request *http.Request) {
	if rt.manager.Config().IsPrintLog() {
		log.Printf("%s %s", request.Method, request.URL.Path)
	}
}

// ValidateRootUrl Checks if the root url matches the "/" character.
func ValidateRootUrl(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return false
	}
	return true
}

// ParseSlugIndex parses url fragments and records their indexes as keys in the map.
// The values are bool values: if true - this fragment is a slug field, if false - it is a regular fragment.
func ParseSlugIndex(path []string) map[int]bool {
	res := make(map[int]bool)
	for i := 0; i < len(path); i++ {
		if len(path) < 2 {
			continue
		}
		if string(path[i][0]) == "<" && string(path[i][len(path[i])-1]) == ">" {
			res[i] = true
		} else {
			res[i] = false
		}
	}
	return res
}

// HandleSlugUrls by number, inserts values from the current addressee into the url pattern in place of slug values.
// Also, sets slug parameters as a map.
func HandleSlugUrls(parseUrl map[int]bool, slugUrl []string, url []string) (string, map[string]string) {
	if len(slugUrl) != len(url) {
		return "", nil
	}
	if len(slugUrl) != len(parseUrl) {
		return "", nil
	}
	params := make(map[string]string)
	for i, isSlug := range parseUrl {
		if isSlug {
			params[strings.Trim(slugUrl[i], "<>")] = url[i]
			slugUrl[i] = url[i]
		}
	}
	res := "/" + strings.Join(slugUrl, "/")
	return res, params
}

// SplitUrlFromFirstSlug returns the left side of the url before the "<" sign.
func splitUrlFromFirstSlug(url string) string {
	index := strings.Index(url, "<")
	if index == -1 {
		return url
	}
	return url[:index]
}
