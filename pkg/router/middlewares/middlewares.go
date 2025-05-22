package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"sync"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

var mu sync.Mutex

type PreMiddleware func(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error
type PostMiddleware func(r *http.Request, m interfaces.IManager) error
type AsyncMiddleware func(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error

type IMiddleware interface {
	PreMiddleware(order int, handler PreMiddleware)
	PostMiddleware(order int, handler PostMiddleware)
	AsyncMiddleware(handler AsyncMiddleware)
	RunPreMiddlewares(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error
	RunAndWaitAsyncMiddlewares(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error
	RunPostMiddlewares(r *http.Request, m interfaces.IManager) error
}

// Middlewares implementation of the middleware concept for the framework.
// Three types of middleware are possible:
//  1. PreMiddleware — executed synchronously before the request.
//     These middlewares are executed exactly according to their established order.
//  2. PostMiddleware — executed synchronously after the request. These middlewares
//     will be called in the specified order.
//  3. AsyncMiddleware — asynchronous middleware. They are executed asynchronously,
//     but after [PreMiddleware] and before the request handler. Cannot be ordered, so are not used in a chain.
//
// Also each middleware returns an error, it is handled in the router. But if you need more specific processing,
// you can not return an error, but process it directly in the middleware.
//
// Also very important point: there are also functions that can control the router from middleware,
// they are in the same package. These functions should be called only in middleware.
// These functions include:
//  1. SkipNextPage — skips the page turner. That is, not just its display, but the whole logic.
//  2. SkipNextPageAndRedirect — does the same thing as [SkipNextPage], but does redirect after skipping the page.
type Middlewares struct {
	preMiddlewaresOrder  sort.IntSlice
	postMiddlewaresOrder sort.IntSlice
	preMiddlewares       map[int]PreMiddleware
	postMiddlewares      map[int]PostMiddleware
	asyncMiddlewares     []AsyncMiddleware
}

func NewMiddlewares() *Middlewares {
	return &Middlewares{
		preMiddlewares:  make(map[int]PreMiddleware),
		postMiddlewares: make(map[int]PostMiddleware),
	}
}

// PreMiddleware handlers that are executed in an ordered fashion before the url handler.
// The order must not be repeated.
func (mddl *Middlewares) PreMiddleware(order int, handler PreMiddleware) {
	if slices.Contains(mddl.preMiddlewaresOrder, order) {
		panic(fmt.Sprintf("middleware with order %s already exists", strconv.Itoa(order)))
	}
	mddl.preMiddlewaresOrder = append(mddl.preMiddlewaresOrder, order)
	mddl.preMiddlewares[order] = handler
}

// PostMiddleware processing is performed after url processing is finished.
func (mddl *Middlewares) PostMiddleware(order int, handler PostMiddleware) {
	if slices.Contains(mddl.postMiddlewaresOrder, order) {
		panic(fmt.Sprintf("middleware with order %s already exists", strconv.Itoa(order)))
	}
	mddl.postMiddlewaresOrder = append(mddl.postMiddlewaresOrder, order)
	mddl.postMiddlewares[order] = handler
}

// AsyncMiddleware handler that is executed asynchronously before the request handler,
// but after [PreMiddleware] processing.
// Can't create chains, not called in an orderly fashion.
func (mddl *Middlewares) AsyncMiddleware(handler AsyncMiddleware) {
	mddl.asyncMiddlewares = append(mddl.asyncMiddlewares, handler)
}

func (mddl *Middlewares) RunPreMiddlewares(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error {
	mddl.preMiddlewaresOrder.Sort()
	for i := 0; i < mddl.preMiddlewaresOrder.Len(); i++ {
		order := mddl.preMiddlewaresOrder[i]
		if err := mddl.preMiddlewares[order](w, r, m); err != nil {
			return err
		}
	}
	return nil
}

// RunAndWaitAsyncMiddlewares runs asynchronous middlewares.
// It also waits for them to complete, no additional actions are needed.
func (mddl *Middlewares) RunAndWaitAsyncMiddlewares(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error {
	var wg sync.WaitGroup
	var asyncError error
	stopAsyncMiddlewares := make(chan struct{})
	for i := 0; i < len(mddl.asyncMiddlewares); i++ {
		handler := mddl.asyncMiddlewares[i]
		wg.Add(1)
		go func(h AsyncMiddleware) {
			defer wg.Done()

			// If at least one handler causes an error, all other handlers will fail to run.
			select {
			case <-stopAsyncMiddlewares:
				return
			default:
			}

			if err := h(w, r, m); err != nil {
				mu.Lock()
				asyncError = err
				close(stopAsyncMiddlewares)
				mu.Unlock()
			}
		}(handler)
	}
	wg.Wait()
	return asyncError
}

func (mddl *Middlewares) RunPostMiddlewares(r *http.Request, m interfaces.IManager) error {
	mddl.postMiddlewaresOrder.Sort()
	for i := 0; i < mddl.postMiddlewaresOrder.Len(); i++ {
		order := mddl.postMiddlewaresOrder[i]
		if err := mddl.postMiddlewares[order](r, m); err != nil {
			return err
		}
	}
	return nil
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
