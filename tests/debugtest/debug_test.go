package debugtest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
)

var mngr = manager.NewManager(nil)

func TestMain(m *testing.M) {
	mngr.Config().ErrorLoggingFile("test.log")
	newRouter := router.NewRouter(mngr)
	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.Get("/server-err", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerError(w, "error", manager.Config()) }
	})
	newRouter.Get("/server-forbidden", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerForbidden(w, manager.Config()) }
	})
	newRouter.Get("/server-logging", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerError(w, "Logging test", manager.Config()) }
	})
	serv := server.NewServer(":8040", newRouter)
	go func() {
		err := serv.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(":8040", 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err := serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestServerErrorDebTrue(t *testing.T) {
	mngr.Config().Debug(true)
	get, err := http.Get("http://localhost:8040/server-err")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "error" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestServerErrorDebFalse(t *testing.T) {
	mngr.Config().Debug(false)
	get, err := http.Get("http://localhost:8040/server-err")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "500 Internal server error" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestServerForbidden(t *testing.T) {
	mngr.Config().Debug(false)
	get, err := http.Get("http://localhost:8040/server-forbidden")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "403 forbidden" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestLogging(t *testing.T) {
	mngr.Config().ErrorLogging(true)
	mngr.Config().Debug(true)
	get, err := http.Get("http://localhost:8040/server-logging")
	if err != nil {
		t.Error(err)
	}
	if fstring.PathExist("test.log") {
		os.Remove("test.log")
	} else {
		t.Errorf("Log file not found.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
	mngr.Config().ErrorLogging(false)
}
