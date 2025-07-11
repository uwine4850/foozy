package routersimplehandlers_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
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
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.HandlerSet(map[string][]map[string]router.Handler{
		router.MethodGET: {
			{"/handler-set-get-1": func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
				w.Write([]byte("OK1"))
				return nil
			}},
			{"/handler-set-get-2": func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
				w.Write([]byte("OK2"))
				return nil
			}},
		},
		router.MethodDELETE: {
			{"/handler-set-delete-1": func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
				w.Write([]byte("OK1"))
				return nil
			}},
			{"/handler-set-delete-2": func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
				w.Write([]byte("OK2"))
				return nil
			}},
		},
	})
	newRouter.Register(router.MethodGET, "/test-get", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodPOST, "/test-post", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodDELETE, "/test-delete", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodPUT, "/test-put", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodPATCH, "/test-patch", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodOPTIONS, "/test-options", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodHEAD, "/test-head", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-error", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		return errors.New("ERROR")
	})
	newRouter.Register(router.MethodGET, "/test-slug/:slug", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		slug, _ := manager.OneTimeData().GetSlugParams("slug")
		w.Write([]byte(slug))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-multiple-slug/:slug/:slug1", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		slug, _ := manager.OneTimeData().GetSlugParams("slug")
		slug1, _ := manager.OneTimeData().GetSlugParams("slug1")
		w.Write([]byte(slug + " " + slug1))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-redirect-error", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		router.RedirectError(w, r, "/catch-redirect-error", "test error")
		return nil
	})
	newRouter.Register(router.MethodGET, "/catch-redirect-error", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		router.CatchRedirectError(r, manager)
		errorText, _ := manager.OneTimeData().GetUserContext(namelib.ROUTER.REDIRECT_ERROR)
		w.Write([]byte(errorText.(string)))
		return nil
	})
	newServer := server.NewServer(tutils.PortSimpleHandlers, newRouter, nil)
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

func TestGet(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortSimpleHandlers, "test-get"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("GET request processing error")
	}
}

func TestPost(t *testing.T) {
	resp, err := form.SendApplicationForm(tutils.MakeUrl(tutils.PortSimpleHandlers, "test-post"), map[string][]string{})
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("POST request processing error")
	}
}

func TestDelete(t *testing.T) {
	resp, err := tutils.SendRequest(http.MethodDelete, tutils.MakeUrl(tutils.PortSimpleHandlers, "test-delete"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("DELETE request processing error")
	}
}

func TestPut(t *testing.T) {
	resp, err := tutils.SendRequest(http.MethodPut, tutils.MakeUrl(tutils.PortSimpleHandlers, "test-put"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("PUT request processing error")
	}
}

func TestPatch(t *testing.T) {
	resp, err := tutils.SendRequest(http.MethodPatch, tutils.MakeUrl(tutils.PortSimpleHandlers, "test-patch"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("PATCH request processing error")
	}
}

func TestOptions(t *testing.T) {
	resp, err := tutils.SendRequest(http.MethodOptions, tutils.MakeUrl(tutils.PortSimpleHandlers, "test-options"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("OPTIONS request processing error")
	}
}

func TestHead(t *testing.T) {
	resp, err := tutils.SendRequest(http.MethodHead, tutils.MakeUrl(tutils.PortSimpleHandlers, "test-head"))
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("HEAD request processing error")
	}
}

func TestError(t *testing.T) {
	resp, err := tutils.SendRequest(http.MethodGet, tutils.MakeUrl(tutils.PortSimpleHandlers, "test-error"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("status code must be 500")
	}
	if res != "ERROR" {
		t.Errorf("the error text does not match the expected text")
	}
}

func TestHandlerSet(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		expect string
	}{
		{"handler test GET 1", http.MethodGet, tutils.MakeUrl(tutils.PortSimpleHandlers, "handler-set-get-1"), "OK1"},
		{"handler test GET 2", http.MethodGet, tutils.MakeUrl(tutils.PortSimpleHandlers, "handler-set-get-2"), "OK2"},
		{"handler test DELETE 1", http.MethodGet, tutils.MakeUrl(tutils.PortSimpleHandlers, "handler-set-get-2"), "OK2"},
		{"handler test DELETE 2", http.MethodGet, tutils.MakeUrl(tutils.PortSimpleHandlers, "handler-set-get-2"), "OK2"},
	}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			resp, err := tutils.SendRequest(test.method, test.url)
			if err != nil {
				t.Error(err)
			}
			res, err := tutils.ReadBody(resp.Body)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != 200 {
				t.Errorf("status code must be 200")
			}
			if res != test.expect {
				t.Errorf("got %s want %s", res, test.expect)
			}
		})
	}
}

func TestSlug(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortSimpleHandlers, "test-slug/slug"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "slug" {
		t.Errorf("slug does not match the expectation")
	}
}

func TestMultipleSlug(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortSimpleHandlers, "test-multiple-slug/slug/slug1"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "slug slug1" {
		t.Errorf("slug does not match the expectation")
	}
}

func TestRedirectError(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortSimpleHandlers, "test-redirect-error"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "test error" {
		t.Errorf("error does not match the expectation")
	}
}
