package csrf_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/secure"
	"github.com/uwine4850/foozy/pkg/server"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

var cookieToken string

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		nil,
		database.NewDatabasePool(),
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/set-cookie-token", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		token, err := secure.GenerateCsrfToken()
		if err != nil {
			return err
		}
		cookieToken = token
		cookie := &http.Cookie{
			Name:     namelib.ROUTER.COOKIE_CSRF_TOKEN,
			Value:    token,
			MaxAge:   1000,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-validate-cookie-token", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		if err := secure.ValidateCookieCsrfToken(r, cookieToken); err != nil {
			return err
		}
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-set-csrf-token", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		if err := secure.SetCSRFToken(1000, true, w, r, manager); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/read-token", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		_, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
		if err != nil {
			return err
		}
		w.Write([]byte("OK"))
		return nil
	})
	newServer := server.NewServer(tutils.PortCSRFToken, newRouter, nil)
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

func TestGenerateCsrfToken(t *testing.T) {
	token, err := secure.GenerateCsrfToken()
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("csrf token not generated")
	}
}

func TestValidateCookieToken(t *testing.T) {
	writeCoolieResponse, err := http.Get(tutils.MakeUrl(tutils.PortCSRFToken, "set-cookie-token"))
	if err != nil {
		t.Error(err)
	}
	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tutils.PortCSRFToken, "test-validate-cookie-token"), nil)
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
		t.Errorf("token validation error: %s", res)
	}
}

func TestSetCSRFToken(t *testing.T) {
	writeCoolieResponse, err := http.Get(tutils.MakeUrl(tutils.PortCSRFToken, "test-set-csrf-token"))
	if err != nil {
		t.Error(err)
	}
	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tutils.PortCSRFToken, "read-token"), nil)
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
		t.Errorf("token not found: %s", res)
	}
}
