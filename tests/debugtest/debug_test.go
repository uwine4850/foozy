package debugtest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
	initcnf "github.com/uwine4850/foozy/tests/init_cnf"
)

var mngr = manager.NewManager(nil)

func TestMain(m *testing.M) {
	initcnf.InitCnf()
	newRouter := router.NewRouter(mngr)
	newRouter.Get("/server-forbidden", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerForbidden(w, manager) }
	})
	newRouter.Get("/server-logging", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerError(w, "Logging test", manager) }
	})
	serv := server.NewServer(":8040", newRouter, nil)
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

func TestServerForbidden(t *testing.T) {
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
	get, err := http.Get("http://localhost:8040/server-logging")
	if err != nil {
		t.Error(err)
	}
	if fpath.PathExist(config.LoadedConfig().Default.Debug.ErrorLoggingPath) {
		os.Remove(config.LoadedConfig().Default.Debug.ErrorLoggingPath)
	} else {
		t.Errorf("Log file not found.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
