package router

import (
	"fmt"
	"github.com/uwine4850/foozy/internal/interfaces"
	"github.com/uwine4850/foozy/internal/tmlengine"
	"github.com/uwine4850/foozy/internal/utils"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	mux            http.ServeMux
	request        *http.Request
	writer         http.ResponseWriter
	TemplateEngine tmlengine.ITemplateEngine
	templatePath   string
	context        map[string]interface{}
	manager        interfaces.IManager
	enableLog      bool
	middleware     interfaces.IMiddleware
}

func NewRouter(manager interfaces.IManager) *Router {
	return &Router{mux: *http.NewServeMux(), manager: manager}
}

// Get Processing a GET request. Called only once.
func (rt *Router) Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "GET", fn))
}

// Post Processing a POST request. Called only once.
func (rt *Router) Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "POST", fn))
}

func (rt *Router) GetMux() *http.ServeMux {
	return &rt.mux
}

// getHandleFunc This method handles each http method call.
func (rt *Router) getHandleFunc(pattern string, method string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rt.setWR(writer, request)
		if !rt.validateMethod(method) {
			return
		}
		rt.printLog(request)

		// Check if the url matches its pattern with possible slug fields.
		if request.URL.Path != "/" {
			parseUrl := ParseSlugIndex(utils.SplitUrl(pattern))
			res, params := HandleSlugUrls(parseUrl, utils.SplitUrl(pattern), utils.SplitUrl(request.URL.Path))
			if res != request.URL.Path {
				http.NotFound(writer, request)
				return
			}
			if params != nil {
				rt.manager.SetSlugParams(params)
			}
		}
		if rt.middleware != nil {
			err := rt.middleware.RunPreMddl(writer, request, rt.manager)
			if err != nil {
				panic(err)
			}
		}
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

func (rt *Router) SetMiddleware(middleware interfaces.IMiddleware) {
	rt.middleware = middleware
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
