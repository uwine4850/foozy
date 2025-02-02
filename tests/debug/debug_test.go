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
	"github.com/uwine4850/foozy/tests1/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests1/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

var mngr = manager.NewManager(nil)

func TestMain(m *testing.M) {
	testinitcnf.InitCnf()
	newRouter := router.NewRouter(mngr)
	newRouter.Get("/server-forbidden", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerForbidden(w, manager) }
	})
	newRouter.Get("/server-logging", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { router.ServerError(w, "Logging test", manager) }
	})
	serv := server.NewServer(tconf.PortDebug, newRouter, nil)
	go func() {
		err := serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortDebug, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
	err := serv.Stop()
	if err != nil {
		panic(err)
	}
}

func TestServerForbidden(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortDebug, "server-forbidden"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortDebug, "server-logging"))
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
