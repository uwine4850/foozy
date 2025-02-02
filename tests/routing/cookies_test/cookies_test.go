package cookiestest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/tests1/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests1/common/test_init_cnf"
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
	testinitcnf.InitCnf()
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newRouter := router.NewRouter(manager.NewManager(render))
	newRouter.Get("/session-create", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		if err := cookies.CreateSecureCookieData(hashKey, blockKey, w, &http.Cookie{
			Name:     "session",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}, &SessionData{UserID: "111"}); err != nil {
			panic(err)
		}
		return func() {}
	})
	newRouter.Get("/session-read", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		var data SessionData
		if err := cookies.ReadSecureCookieData(hashKey, blockKey, r, "session", &data); err != nil {
			panic(err)
		}
		w.Write([]byte(data.UserID))
		return func() {}
	})
	newRouter.Get("/cookie", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { cookies.SetStandartCookie(w, "cookie", "value", "/", 0) }
	})
	serv := server.NewServer(tconf.PortCookies, newRouter, nil)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortCookies, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestLoginSession(t *testing.T) {
	createReq, err := http.NewRequest("GET", tutils.MakeUrl(tconf.PortCookies, "session-create"), nil)
	if err != nil {
		t.Error(err)
	}
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Error(err)
	}
	defer createResp.Body.Close()

	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tconf.PortCookies, "session-read"), nil)
	if err != nil {
		t.Error(err)
	}

	readReq.AddCookie(createResp.Cookies()[0])

	readResp, err := http.DefaultClient.Do(readReq)
	if err != nil {
		t.Error(err)
	}
	defer readResp.Body.Close()

	body, err := io.ReadAll(readResp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "111" {
		t.Errorf("The secure session data was not read correctly.")
	}
}

func TestSetStandartCookie(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", tutils.MakeUrl(tconf.PortCookies, "cookie"), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookies to be set, but none were found")
	}

	var mycookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "cookie" {
			mycookie = cookie
			break
		}
	}

	if mycookie == nil {
		t.Fatal("Expected mycookie to be set, but it was not found")
	}

	if mycookie.Value != "value" {
		t.Errorf("Expected mycookie to have value %v, but got %v", "value", mycookie.Value)
	}
}
