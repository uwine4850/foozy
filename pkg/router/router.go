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
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if config.LoadedConfig().Default.Debug.Debug {
		debug.ErrorLoggingIfEnableAndWrite(w, err.Error(), err.Error())
	} else {
		debug.ErrorLoggingIfEnableAndWrite(w, err.Error(), "500 Internal server error")
	}
}

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

var managerObject interfaces.IManager = nil

type Router struct {
	mux               http.ServeMux
	routes            map[string]muxRouter
	middleware        interfaces.IMiddleware
	internalErrorFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewRouter(manager interfaces.IManager) *Router {
	managerObject = manager
	return &Router{mux: *http.NewServeMux(), routes: map[string]muxRouter{}, internalErrorFunc: internalServerError}
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

// InternalError sets the function to be used when handling internal errors.
func (rt *Router) InternalError(fn func(w http.ResponseWriter, r *http.Request, err error)) {
	rt.internalErrorFunc = fn
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
		if err := debug.ClearRequestInfoLogging(); err != nil {
			rt.internalErrorFunc(writer, request, err)
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, fmt.Sprintf("request url: %s", request.URL))
		manager, err := rt.initManager()
		if err != nil {
			rt.internalErrorFunc(writer, request, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "manager is initialized")

		manager.OneTimeData().SetUserContext(namelib.ROUTER.URL_PATTERN, urlPattern)

		parseUrl := ParseSlugIndex(fstring.SplitUrl(urlPattern))
		if request.URL.Path != "/" && len(parseUrl) > 0 {
			res, params := HandleSlugUrls(parseUrl, fstring.SplitUrl(urlPattern), fstring.SplitUrl(request.URL.Path))
			if res != request.URL.Path {
				http.NotFound(writer, request)
				return
			}
			if params != nil {
				manager.OneTimeData().SetSlugParams(params)
			}
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "slug url is parsed")

		// Run middlewares.
		if skip, err := rt.runMddl(writer, request, manager); err != nil {
			if rt.internalErrorFunc != nil {
				rt.internalErrorFunc(writer, request, err)
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			} else {
				rt.internalErrorFunc(writer, request, err)
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			}
			return
		} else {
			if skip {
				return
			}
		}
		rt.switchRegisterMethods(writer, request, _muxRouter, manager)
		rt.printLog(request)
	}
}

func (rt *Router) switchRegisterMethods(writer http.ResponseWriter, request *http.Request, _muxRouter muxRouter, manager interfaces.IManager) {
	debug.RequestLogginIfEnable(debug.P_ROUTER, "run switch methods...")
	debug.RequestLogginIfEnable(debug.P_ROUTER, "validate method")
	connection := request.Header.Get("Connection")
	if connection != "" && connection == "Upgrade" {
		handler := _muxRouter.Ws
		if !rt.validateMethod(handler, "WS", writer) {
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run WS handler")
		handler(writer, request, manager)()
		return
	}
	switch request.Method {
	case http.MethodGet:
		handler := _muxRouter.Get
		if !rt.validateMethod(handler, "GET", writer) {
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run GET handler")
		handler(writer, request, manager)()
	case http.MethodPost:
		handler := _muxRouter.Post
		if !rt.validateMethod(handler, "POST", writer) {
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run POST handler")
		handler(writer, request, manager)()
	case http.MethodPut:
		handler := _muxRouter.Put
		if !rt.validateMethod(handler, "PUT", writer) {
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run PUT handler")
		handler(writer, request, manager)()
	case http.MethodDelete:
		handler := _muxRouter.Delete
		if !rt.validateMethod(handler, "DELETE", writer) {
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run DELETE handler")
		handler(writer, request, manager)()
	case http.MethodOptions:
		handler := _muxRouter.Options
		if !rt.validateMethod(handler, "OPTIONS", writer) {
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run OPTIONS handler")
		handler(writer, request, manager)()
	}
	debug.RequestLogginIfEnable(debug.P_ROUTER, "method completed")
}

// initManager initializes a new manager instance.
// Must be called for each new request.
func (rt *Router) initManager() (interfaces.IManager, error) {
	debug.RequestLogginIfEnable(debug.P_ROUTER, "init manager")
	_newManager, err := managerObject.New()
	if err != nil {
		return nil, err
	}
	newManager := _newManager.(interfaces.IManager)
	// Set OneTimeData.
	newOTD, err := manager.CreateNewManagerData(managerObject)
	if err != nil {
		return nil, err
	}
	// Set render.
	newManager.SetOneTimeData(newOTD)
	if managerObject.Render() != nil {
		newRender, err := tmlengine.CreateNewRenderInstance(managerObject)
		if err != nil {
			return nil, err
		}
		newManager.SetRender(newRender)
	}
	return newManager, nil
}

// runMddl runs the middleware.
// Conventional middleware is run first because it is more prioritized than asynchronous. They usually
// perform important logic, which, for example, should be run first.
// After execution of synchronous middleware, asynchronous ones are executed.
// Middleware errors and the page rendering skip algorithm are also handled here.
func (rt *Router) runMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, error) {
	if rt.middleware != nil {
		debug.RequestLogginIfEnable(debug.P_ROUTER, "run middlewares...")
		// Running synchronous middleware.
		err := rt.middleware.RunMddl(w, r, manager)
		if err != nil {
			return false, err
		}
		// Running asynchronous middleware.
		rt.middleware.RunAsyncMddl(w, r, manager)
		// Waiting for all asynchronous middleware to complete.
		rt.middleware.WaitAsyncMddl()
		// Handling middleware errors.
		mddlErr, err := middlewares.GetMddlError(manager.OneTimeData())
		if err != nil {
			return false, err
		}
		if mddlErr != nil {
			return false, errors.New(mddlErr.Error())
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "middlewares are completed")
		// Checking the skip of the next page. Runs after a more important error check.
		if middlewares.IsSkipNextPage(manager.OneTimeData()) {
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
	if config.LoadedConfig().Default.Debug.PrintInfo {
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
