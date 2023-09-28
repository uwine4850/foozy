package router

import (
	"fmt"
	"github.com/uwine4850/foozy/internal/tmlengine"
	"log"
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
	enableLog      bool
}

func NewRouter(manager IManager) *Router {
	return &Router{mux: *http.NewServeMux(), manager: manager}
}

// Get Processing a GET request. Called only once.
func (rt *Router) Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager)) {
	rt.pattern = pattern
	rt.mux.Handle(pattern, rt.getHandleFunc("GET", fn))
}

// Post Processing a POST request. Called only once.
func (rt *Router) Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager)) {
	rt.pattern = pattern
	rt.mux.Handle(pattern, rt.getHandleFunc("POST", fn))
}

func (rt *Router) GetMux() *http.ServeMux {
	return &rt.mux
}

// getHandleFunc This method handles each http method call.
func (rt *Router) getHandleFunc(method string, fn func(w http.ResponseWriter, r *http.Request, manager IManager)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rt.setWR(writer, request)
		if !rt.validateMethod(method) {
			return
		}
		rt.printLog(request)
		fn(writer, request, rt.manager)
	}
}

// validateMethod Check if the http method matches the expected method.
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

func (rt *Router) EnableLog(enable bool) {
	rt.enableLog = enable
}

func (rt *Router) printLog(request *http.Request) {
	if rt.enableLog {
		log.Println(fmt.Sprintf("%s %s", request.Method, request.URL.Path))
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
