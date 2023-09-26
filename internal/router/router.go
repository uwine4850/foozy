package router

import (
	"fmt"
	"net/http"
)

type IRouter interface {
	Get(pattern string, fn func(w http.ResponseWriter, r *http.Request))
	Post(pattern string, fn func(w http.ResponseWriter, r *http.Request))
	GetMux() *http.ServeMux
}

type Router struct {
	mux     http.ServeMux
	request *http.Request
	writer  http.ResponseWriter
	pattern string
}

func NewRouter() *Router {
	return &Router{mux: *http.NewServeMux()}
}

func (rt *Router) Get(pattern string, fn func(w http.ResponseWriter, r *http.Request)) {
	rt.pattern = pattern
	rt.mux.Handle(pattern, rt.getHandleFunc("GET", fn))
}

func (rt *Router) Post(pattern string, fn func(w http.ResponseWriter, r *http.Request)) {
	rt.pattern = pattern
	rt.mux.Handle(pattern, rt.getHandleFunc("POST", fn))
}

func (rt *Router) GetMux() *http.ServeMux {
	return &rt.mux
}

func (rt *Router) getHandleFunc(method string, fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rt.setWR(writer, request)
		if !rt.validateMethod(method) {
			return
		}
		if !rt.validateRootUrl() {
			return
		}
		fn(writer, request)
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
