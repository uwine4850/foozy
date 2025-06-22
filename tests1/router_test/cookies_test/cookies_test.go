package cookies_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/server"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

type SessionData struct {
	UserID string
}

var (
	hashKey  = []byte("1234567890abcdef1234567890abcdef") // 32 bytes
	blockKey = []byte("abcdefghijklmnopqrstuvwx12345678") // 32 bytes
)

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		nil,
		manager.NewDatabasePool(),
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/test-write-secure-cookie", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		if err := cookies.CreateSecureCookieData(hashKey, blockKey, w, &http.Cookie{
			Name:     "session",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}, &SessionData{UserID: "111"}); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-read-secure-cookie", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		var data SessionData
		if err := cookies.ReadSecureCookieData(hashKey, blockKey, r, "session", &data); err != nil {
			return err
		}
		w.Write([]byte(data.UserID))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-nohmac-write-secure-cookie", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		err := cookies.CreateSecureNoHMACCookieData(blockKey, w, &http.Cookie{
			Name:     "session",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}, &SessionData{UserID: "111"})
		if err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-nohmac-read-secure-cookie", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		var data SessionData
		if err := cookies.ReadSecureNoHMACCookieData(blockKey, r, "session", &data); err != nil {
			return err
		}
		w.Write([]byte(data.UserID))
		return nil
	})
	newServer := server.NewServer(tutils.PortCookies, newRouter, nil)
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
	writeCoolieResponse, err := http.Get(tutils.MakeUrl(tutils.PortCookies, "test-write-secure-cookie"))
	if err != nil {
		t.Error(err)
	}
	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tutils.PortCookies, "test-read-secure-cookie"), nil)
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
	if res != "111" {
		t.Error("error reading or writing a secure cookie ", res)
	}
}

func TestNoHMACSecureData(t *testing.T) {
	writeCoolieResponse, err := http.Get(tutils.MakeUrl(tutils.PortCookies, "test-nohmac-write-secure-cookie"))
	if err != nil {
		t.Error(err)
	}
	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tutils.PortCookies, "test-nohmac-read-secure-cookie"), nil)
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
	if res != "111" {
		t.Error("error reading or writing a secure cookie ", res)
	}
}
