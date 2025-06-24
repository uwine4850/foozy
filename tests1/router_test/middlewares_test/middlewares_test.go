package middlewares_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/server"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

var preMddl = false
var asyncMddl = false
var postMddl = false
var newMiddlewares = middlewares.NewMiddlewares()

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		nil,
		manager.NewDatabasePool(),
	)
	newMiddlewares.PreMiddleware(0, func(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error {
		preMddl = true
		return nil
	})
	newMiddlewares.AsyncMiddleware(func(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error {
		asyncMddl = true
		return nil
	})
	newMiddlewares.PostMiddleware(0, func(r *http.Request, m interfaces.IManager) error {
		postMddl = true
		return nil
	})
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/test-run", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		if preMddl && asyncMddl {
			w.Write([]byte("OK"))
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-skip-next-page", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		w.Write([]byte("PAGE"))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-skip-next-page-redirect", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		w.Write([]byte("NOT DISPLAY PAGE"))
		return nil
	})
	newRouter.Register(router.MethodGET, "/redirect", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		w.Write([]byte("REDIRECT PAGE"))
		return nil
	})
	newServer := server.NewServer(tutils.PortMiddlewares, newRouter, nil)
	go func() {
		if err := newServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	exitCode := m.Run()
	if err := newServer.Stop(); err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func onError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func TestRunMiddlewares(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortMiddlewares, "test-run"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" && postMddl {
		t.Errorf("middlewares run error")
	}
}

func TestSkipNextPage(t *testing.T) {
	newMiddlewares.AsyncMiddleware(func(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error {
		middlewares.SkipNextPage(m.OneTimeData())
		return nil
	})
	resp, err := http.Get(tutils.MakeUrl(tutils.PortMiddlewares, "test-skip-next-page"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res == "PAGE" {
		t.Error("the page showed up, but it shouldn't have")
	}
}

func TestSkipNextPageAndRedirect(t *testing.T) {
	newMiddlewares.AsyncMiddleware(func(w http.ResponseWriter, r *http.Request, m interfaces.IManager) error {
		if r.URL.Path == "/test-skip-next-page-redirect" {
			middlewares.SkipNextPageAndRedirect(m.OneTimeData(), w, r, "/redirect")
		}
		return nil
	})
	resp, err := http.Get(tutils.MakeUrl(tutils.PortMiddlewares, "test-skip-next-page-redirect"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "REDIRECT PAGE" {
		t.Error("page is not redirected")
	}
}
