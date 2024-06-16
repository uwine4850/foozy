package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	interfaces2 "github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/utils"
)

type Router struct {
	mux            http.ServeMux
	request        *http.Request
	writer         http.ResponseWriter
	TemplateEngine interfaces2.ITemplateEngine
	manager        interfaces2.IManager
	enableLog      bool
	middleware     interfaces2.IMiddleware
}

func NewRouter(manager interfaces2.IManager) *Router {
	return &Router{mux: *http.NewServeMux(), manager: manager}
}

// Get Processing a GET request. Called only once.
func (rt *Router) Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "GET", nil, fn))
}

// Post Processing a POST request. Called only once.
func (rt *Router) Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "POST", nil, fn))
}

func (rt *Router) Put(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "PUT", nil, fn))
}

func (rt *Router) Delete(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "DELETE", nil, fn))
}

// Ws Processing a websocket connection. Used only for communication with the client's websocket.
func (rt *Router) Ws(pattern string, ws interfaces2.IWebsocket, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()) {
	rt.mux.Handle(utils.SplitUrlFromFirstSlug(pattern), rt.getHandleFunc(pattern, "WS", ws, fn))
}

func (rt *Router) GetMux() *http.ServeMux {
	return &rt.mux
}

// getHandleFunc This method handles each http method call.
func (rt *Router) getHandleFunc(pattern string, method string, ws interfaces2.IWebsocket, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager) func()) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rt.setWR(writer, request)
		if !rt.validateMethod(method) {
			return
		}
		rt.printLog(request)

		rt.manager.SetOneTimeData(manager.NewManagerData())

		rt.manager.OneTimeData().SetUserContext("URL_PATTERN", pattern)

		// Check if the url matches its pattern with possible slug fields.
		parseUrl := ParseSlugIndex(utils.SplitUrl(pattern))
		if request.URL.Path != "/" && len(parseUrl) > 0 {
			res, params := HandleSlugUrls(parseUrl, utils.SplitUrl(pattern), utils.SplitUrl(request.URL.Path))
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
			ServerError(writer, err.Error(), rt.manager.Config())
			return
		} else {
			if skip {
				return
			}
		}
		if method == "WS" {
			rt.manager.WS().SetWebsocket(ws)
		}
		mustCall := fn(writer, request, rt.manager)
		mustCall()
	}
}

// validateMethod Check if the http method matches the expected method.
func (rt *Router) validateMethod(method string) bool {
	if method == "WS" {
		return true
	}
	if rt.request.Method != method {
		rt.writer.Header().Set("Allow", method)
		http.Error(rt.writer, fmt.Sprintf("Method %s Not Allowed", rt.request.Method), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (rt *Router) setWR(w http.ResponseWriter, r *http.Request) {
	rt.writer = w
	rt.request = r
}

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
			rt.manager.OneTimeData().DelUserContext("mddlerr")
			return false, errors.New(mddlErr.Error())
		}
		// Checking the skip of the next page. Runs after a more important error check.
		if middlewares.IsSkipNextPage(rt.manager.OneTimeData()) {
			rt.manager.OneTimeData().DelUserContext("skipNextPage")
			return true, nil
		}
	}
	return false, nil
}

// SetTemplateEngine sets the template engine interface.
func (rt *Router) SetTemplateEngine(engine interfaces2.ITemplateEngine) {
	rt.TemplateEngine = engine
}

// EnableLog enable or disable request log output.
func (rt *Router) EnableLog(enable bool) {
	rt.enableLog = enable
}

func (rt *Router) printLog(request *http.Request) {
	if rt.enableLog {
		log.Printf("%s %s", request.Method, request.URL.Path)
	}
}

// SetMiddleware installs the middleware for the handlers.
func (rt *Router) SetMiddleware(middleware interfaces2.IMiddleware) {
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
