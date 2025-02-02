package mddltest

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
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/tests1/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests1/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

var mng = manager.NewManager(nil)
var newRouter = router.NewRouter(mng)

func TestMain(m *testing.M) {
	newRouter.Get("/mddl", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { w.Write([]byte("OK")) }
	})
	newRouter.Get("/redirect", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { w.Write([]byte("is redirect page")) }
	})
	serv := server.NewServer(tconf.PortMddl, newRouter, nil)
	go func() {
		err := serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortMddl, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
	err := serv.Stop()
	if err != nil {
		panic(err)
	}
}

func TestSetGetMddlError(t *testing.T) {
	testinitcnf.InitCnf()
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(0, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		middlewares.SetMddlError(errors.New("mddl error"), manager.OneTimeData())
	})
	newRouter.SetMiddleware(mddl)
	get, err := http.Get(tutils.MakeUrl(tconf.PortMddl, "mddl"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortMddl, "mddl"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortMddl, "mddl"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortMddl, "mddl"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortMddl, "mddl"))
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
