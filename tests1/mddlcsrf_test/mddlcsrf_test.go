package mddlcsrf_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/server"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		nil,
		database.NewDatabasePool(),
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newMiddlewares.PreMiddleware(1, builtin_mddl.GenerateAndSetCsrf(1000, true))
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/test-get-csrf", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		csrfCookie, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
		if err != nil {
			return err
		}
		if csrfCookie.Value == "" {
			w.Write([]byte("csrf token not found"))
		}
		w.Write([]byte("OK"))
		return nil
	})
	newServer := server.NewServer(tutils.PortBuiltinCSRF, newRouter, nil)
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

func TestSecureData(t *testing.T) {
	writeCoolieResponse, err := http.Get(tutils.MakeUrl(tutils.PortBuiltinCSRF, "test-get-csrf"))
	if err != nil {
		t.Error(err)
	}
	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tutils.PortBuiltinCSRF, "test-get-csrf"), nil)
	if err != nil {
		t.Error(err)
	}
	readReq.AddCookie(writeCoolieResponse.Cookies()[0])
	readCookieResponse, err := http.DefaultClient.Do(readReq)
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(readCookieResponse.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Error("error reading or writing a secure cookie ", res)
	}
}
