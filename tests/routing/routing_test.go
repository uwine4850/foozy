package routing

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	router2 "github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	fserer "github.com/uwine4850/foozy/pkg/server"
)

type SessionData struct {
	UserID string
}

var (
	hashKey  = []byte("1234567890abcdef1234567890abcdef") // 32 bytes
	blockKey = []byte("abcdefghijklmnopqrstuvwx12345678") // 32 bytes
)

func TestMain(m *testing.M) {
	newTmplEngine, err := tmlengine.NewTemplateEngine()
	if err != nil {
		panic(err)
	}
	newRouter := router2.NewRouter(manager.NewManager(newTmplEngine))
	newRouter.EnableLog(false)
	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.Get("/page", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		return func() { w.Write([]byte("OK")) }
	})
	newRouter.Get("/page/<id>", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		id, _ := manager.OneTimeData().GetSlugParams("id")
		if id == "1" {
			return func() { w.Write([]byte("OK")) }
		}
		return func() {}
	})
	newRouter.Get("/page2/<id>/<name>", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		id, _ := manager.OneTimeData().GetSlugParams("id")
		name, _ := manager.OneTimeData().GetSlugParams("name")
		if id == "1" && name == "name" {
			return func() { w.Write([]byte("OK")) }
		}
		return func() {}
	})
	newRouter.Post("/post/<id>", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		id, _ := manager.OneTimeData().GetSlugParams("id")
		if id == "12" {
			return func() { w.Write([]byte("OK")) }
		}
		return func() {}
	})
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
	server := fserer.NewServer(":8030", newRouter)
	go func() {
		err = server.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
	exitCode := m.Run()
	err = server.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestPage(t *testing.T) {
	get, err := http.Get("http://localhost:8030/page")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "OK" {
		t.Errorf("Error on page retrieval.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestPageSlug(t *testing.T) {
	get, err := http.Get("http://localhost:8030/page/1")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "OK" {
		t.Errorf("Error receiving slug parameter.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestPageMultipleSlug(t *testing.T) {
	get, err := http.Get("http://localhost:8030/page2/1/name")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "OK" {
		t.Errorf("Error receiving slug parameter.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	resp, err := form.SendApplicationForm("http://localhost:8030/post/12", map[string]string{})
	if err != nil {
		t.Error(err)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "OK" {
		t.Errorf("Error during POST method processing.")
	}
}
