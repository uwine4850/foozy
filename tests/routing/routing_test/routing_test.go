package routingtest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/tests/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests/common/tutils"
)

func TestMain(m *testing.M) {
	testinitcnf.InitCnf()
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newRouter := router.NewRouter(manager.NewManager(render))
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
	newRouter.Get("/redirect-error", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		router.RedirectError(w, r, "/catch-redirect-error", "error")
		return func() {}
	})
	newRouter.Get("/catch-redirect-error", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		router.CatchRedirectError(r, manager)
		_, ok := manager.Render().GetContext()[namelib.ROUTER.REDIRECT_ERROR]
		w.Write([]byte(strconv.FormatBool(ok)))
		return func() {}
	})
	serv := server.NewServer(tconf.PortRouter, newRouter, nil)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortRouter, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestPage(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortRouter, "page"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortRouter, "page/1"))
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
	get, err := http.Get(tutils.MakeUrl(tconf.PortRouter, "page2/1/name"))
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
	resp, err := form.SendApplicationForm(tutils.MakeUrl(tconf.PortRouter, "post/12"), map[string][]string{})
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

func TestRedirectError(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortRouter, "redirect-error"))
	if err != nil {
		t.Error(err)
	}

	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}

	if string(body) == "false" {
		t.Errorf("The error from RedirectError was passed, but was written to the template engine context.")
	}

	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
