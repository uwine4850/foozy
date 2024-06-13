package middlewares_test

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
)

var newTmplEngine, err = tmlengine.NewTemplateEngine()
var mng = manager.NewManager(newTmplEngine)
var newRouter = router.NewRouter(mng)

func TestMain(m *testing.M) {
	if err != nil {
		panic(err)
	}
	mng.Config().Debug(true)
	newRouter.EnableLog(false)
	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.Get("/mddl", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { w.Write([]byte("OK")) }
	})
	newRouter.Get("/redirect", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { w.Write([]byte("is redirect page")) }
	})
	server := server.NewServer(":8050", newRouter)
	go func() {
		err = server.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}()
	exitCode := m.Run()
	err = server.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestSetGetMddlError(t *testing.T) {
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(0, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		middlewares.SetMddlError(errors.New("mddl error"), manager.OneTimeData())
	})
	newRouter.SetMiddleware(mddl)
	get, err := http.Get("http://localhost:8050/mddl")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "mddl error" {
		t.Errorf("Middleware error not handled.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestSkipAndIsSkipNextPage(t *testing.T) {
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(0, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		middlewares.SkipNextPage(manager.OneTimeData())
	})
	newRouter.SetMiddleware(mddl)
	get, err := http.Get("http://localhost:8050/mddl")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "" {
		t.Errorf("The router did not skip rendering the page.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestSkipNextPageAndRedirect(t *testing.T) {
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(0, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		if r.URL.Path == "/mddl" {
			middlewares.SkipNextPageAndRedirect(manager.OneTimeData(), w, r, "/redirect")
		}
	})
	newRouter.SetMiddleware(mddl)
	get, err := http.Get("http://localhost:8050/mddl")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "is redirect page" {
		t.Errorf("Middleware does not redirect the page.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestMiddlewaresIndex(t *testing.T) {
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(3, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	})
	mddl.HandlerMddl(3, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	})
	newRouter.SetMiddleware(mddl)
	get, err := http.Get("http://localhost:8050/mddl")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "Middleware with id 3 already exists." {
		t.Errorf("Id middleware exists, but no error is raised.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestSyncAsyncMiddlewares(t *testing.T) {
	startTime := time.Now()
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(0, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		time.Sleep(200 * time.Millisecond)
	})
	mddl.HandlerMddl(1, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		time.Sleep(200 * time.Millisecond)
	})
	mddl.AsyncHandlerMddl(func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		time.Sleep(200 * time.Millisecond)
	})
	mddl.AsyncHandlerMddl(func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		time.Sleep(200 * time.Millisecond)
	})
	newRouter.SetMiddleware(mddl)
	get, err := http.Get("http://localhost:8050/mddl")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	if elapsedTime > 620*time.Millisecond {
		t.Errorf("Middlewares took too long to execute: %s", elapsedTime.String())
	}

	if string(body) != "OK" {
		t.Errorf("Error while executing middlewares: %s", string(body))
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
