package router

import (
	"bytes"
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

const (
	MethodGET     = "GET"
	MethodPOST    = "POST"
	MethodDELETE  = "DELETE"
	MethodPUT     = "PUT"
	MethodPATCH   = "PATCH"
	MethodHEAD    = "HEAD"
	MethodOPTIONS = "OPTIONS"
)

// BufferedResponseWriter wrapper around [http.ResponseWriter].
// Used to buffer the write into an http response. Now the [Write]
// method doesn't send the response, it just writes it to the buffer.
// To send a response you need to use the [Flush] method.
type BufferedResponseWriter struct {
	original    http.ResponseWriter
	header      http.Header
	statusCode  int
	buffer      bytes.Buffer
	wroteHeader bool
}

func NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		original:   w,
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}
}

func (rw *BufferedResponseWriter) OriginalWriter() http.ResponseWriter {
	return rw.original
}

func (rw *BufferedResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *BufferedResponseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}
	rw.statusCode = statusCode
	rw.wroteHeader = true
}

// Write writes data to the buffer.
func (rw *BufferedResponseWriter) Write(data []byte) (int, error) {
	return rw.buffer.Write(data)
}

// Flush sending the http response of the previously recorded response.
func (rw *BufferedResponseWriter) Flush() (int, error) {
	for k, vv := range rw.header {
		for _, v := range vv {
			rw.original.Header().Add(k, v)
		}
	}
	rw.original.WriteHeader(rw.statusCode)
	return rw.original.Write(rw.buffer.Bytes())
}

// Handler handles the http method.
// Returns an error that is handled by the proper method from [Adapter].
type Handler func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error

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
	manager           interfaces.Manager
	middlewares       middlewares.IMiddleware
	internalErrorFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewAdapter(manager interfaces.Manager, middlewares middlewares.IMiddleware) *Adapter {
	return &Adapter{
		manager:           manager,
		middlewares:       middlewares,
		internalErrorFunc: internalServerError,
	}
}

// Adapt wraps router.Handler in additional functionality.
// It creates a new manager, starts middlewares and does other small operations.
//
// IMPORTANT: if the connection is a websocket connection, [PostMiddlewares] will not work.
// This is done to provide complete security against unexpected behavior.
// In all other cases, [PostMiddlewares] will work as usual.
func (a *Adapter) Adapt(pattern string, handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isWebsocketConn := IsWebsocket(r)
		bw := NewBufferedResponseWriter(w)
		if err := debug.ClearRequestInfoLogging(); err != nil {
			a.internalErrorFunc(bw, r, err)
			a.wrappedFlush(bw, r)
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, fmt.Sprintf("request url: %s", r.URL))
		debug.RequestLogginIfEnable(debug.P_ROUTER, "init manager")
		newManager, err := a.newManager()
		if err != nil {
			a.internalErrorFunc(bw, r, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			a.wrappedFlush(bw, r)
			return
		}
		debug.RequestLogginIfEnable(debug.P_ROUTER, "manager is initialized")

		newManager.OneTimeData().SetUserContext(namelib.ROUTER.URL_PATTERN, pattern)

		// Slug params
		if params := a.getSlugParams(r.URL.Path, pattern); params != nil {
			newManager.OneTimeData().SetSlugParams(params)
		}

		// Run middlewares
		if skip, err := a.runPreAndAsyncMddl(bw, r, newManager); err != nil {
			a.internalErrorFunc(bw, r, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			a.wrappedFlush(bw, r)
			return
		} else if skip {
			a.wrappedFlush(bw, r)
			return
		}

		a.printLog(r)
		if !isWebsocketConn {
			if err := handler(bw, r, newManager); err != nil {
				a.internalErrorFunc(bw, r, err)
			}
			if err := a.middlewares.RunPostMiddlewares(r, newManager); err != nil {
				if !errors.Is(err, middlewares.ErrStopMiddlewares{}) {
					a.internalErrorFunc(bw, r, err)
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					a.wrappedFlush(bw, r)
					return
				}
			}
			if middlewares.IsSkipNextPage(newManager.OneTimeData()) {
				a.wrappedFlush(bw, r)
				return
			} else {
				a.wrappedFlush(bw, r)
			}
		} else {
			if err := handler(bw.OriginalWriter(), r, newManager); err != nil {
				a.internalErrorFunc(bw, r, err)
			}
		}
	}
}

func (a *Adapter) SetOnErrorFunc(fn func(w http.ResponseWriter, r *http.Request, err error)) {
	a.internalErrorFunc = fn
}

func (a *Adapter) newManager() (interfaces.Manager, error) {
	_newManager, err := a.manager.New()
	if err != nil {
		return nil, err
	}
	newManager := _newManager.(interfaces.Manager)
	return newManager, nil
}

// runPreAndAsyncMddl runs the middleware.
// Conventional middleware is run first because it is more prioritized than asynchronous. They usually
// perform important logic, which, for example, should be run first.
// After execution of synchronous middleware, asynchronous ones are executed.
// Middleware errors and the page rendering skip algorithm are also handled here.
func (a *Adapter) runPreAndAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (bool, error) {
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

func (a *Adapter) getSlugParams(currentPath string, pattern string) map[string]string {
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	urlSegments := strings.Split(strings.Trim(currentPath, "/"), "/")
	return MatchUrlSegments(segments, urlSegments)
}

func (a *Adapter) wrappedFlush(bw *BufferedResponseWriter, r *http.Request) {
	if _, err := bw.Flush(); err != nil {
		a.internalErrorFunc(bw.OriginalWriter(), r, err)
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
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

func (r *Router) HandlerSet(handlerSet map[string][]map[string]Handler) {
	for method, handlers := range handlerSet {
		for i := 0; i < len(handlers); i++ {
			for pattern, handler := range handlers[i] {
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

// MatchUrlSegments compares slug segments to the real url.
// If there is a match, it returns a map slug - value.
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

func IsWebsocket(r *http.Request) bool {
	connHdr := strings.ToLower(r.Header.Get("Connection"))
	return strings.Contains(connHdr, "upgrade") &&
		strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}
