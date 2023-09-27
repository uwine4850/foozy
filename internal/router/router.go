package router

import (
	"fmt"
	"github.com/uwine4850/foozy/internal/tmlengine"
	"net/http"
)

type IRouter interface {
	Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	GetMux() *http.ServeMux
	SetTemplateEngine(engine tmlengine.ITemplateEngine)
}

type Router struct {
	mux            http.ServeMux
	request        *http.Request
	writer         http.ResponseWriter
	pattern        string
	TemplateEngine tmlengine.ITemplateEngine
	templatePath   string
	context        map[string]interface{}
	manager        IManager
}

func NewRouter(manager IManager) *Router {
	return &Router{mux: *http.NewServeMux(), manager: manager}
}

func (rt *Router) Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager)) {
	rt.pattern = pattern
	rt.mux.Handle(pattern, rt.getHandleFunc("GET", fn))
}

func (rt *Router) Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager)) {
	rt.pattern = pattern
	rt.mux.Handle(pattern, rt.getHandleFunc("POST", fn))
}

func (rt *Router) GetMux() *http.ServeMux {
	return &rt.mux
}

func (rt *Router) getHandleFunc(method string, fn func(w http.ResponseWriter, r *http.Request, manager IManager)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rt.setWR(writer, request)
		if !rt.validateMethod(method) {
			return
		}
		if !rt.validateRootUrl() {
			return
		}
		fn(writer, request, rt.manager)
	}
}

func (rt *Router) validateRootUrl() bool {
	if rt.pattern == "/" && rt.request.URL.Path != "/" {
		http.NotFound(rt.writer, rt.request)
		return false
	}
	return true
}

func (rt *Router) validateMethod(method string) bool {
	if rt.request.Method != method {
		rt.writer.Header().Set("Allow", method)
		http.Error(rt.writer, fmt.Sprintf("Method %s Not Allowed", rt.request.Method), 405)
		return false
	}
	return true
}

func (rt *Router) setWR(w http.ResponseWriter, r *http.Request) {
	rt.writer = w
	rt.request = r
}

func (rt *Router) SetTemplateEngine(engine tmlengine.ITemplateEngine) {
	rt.TemplateEngine = engine
}
