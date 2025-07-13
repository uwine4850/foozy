package tmplengine_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/flosch/pongo2"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()

	tmlengine.RegisterMultipleGlobalFilter([]tmlengine.Filter{
		{
			Name: "filter1",
			Fn: func(in, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
				var result string
				if in.IsString() {
					result = in.String() + "-filter1"
				}
				return pongo2.AsValue(result), nil
			},
		},
		{
			Name: "filter2",
			Fn: func(in, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
				var result string
				if in.IsString() {
					result = in.String() + "-filter2"
				}
				return pongo2.AsValue(result), nil
			},
		},
	})

	newRender, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		newRender,
		database.NewDatabasePool(),
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/test-render-template", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		manager.Render().SetTemplatePath("test.html")
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-render-no-template", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-template-context", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		manager.Render().SetTemplatePath("test_context.html")
		manager.Render().SetContext(map[string]interface{}{"KEY": "VAL"})
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-render-json", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		if err := manager.Render().RenderJson(map[string]any{"OK": true}, w); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-strslice-filter", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		manager.Render().SetTemplatePath("test_context.html")
		manager.Render().SetContext(map[string]interface{}{"slice": []int{1, 2, 3}})
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-custom-filters", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		manager.Render().SetTemplatePath("test_context.html")
		manager.Render().SetContext(map[string]interface{}{"val1": "value1", "val2": "value2"})
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
	newServer := server.NewServer(tutils.PortPageRender, newRouter, nil)
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

func TestRenderTemplate(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortPageRender, "test-render-template"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("when the page displays an error: %s", res)
	}
}

func TestRenderNoTemplate(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortPageRender, "test-render-no-template"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "The path to the template is not set." {
		t.Error("error is not displayed")
	}
}

func TestTemplateContext(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortPageRender, "test-template-context"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "VAL|" {
		if res == "" {
			t.Error("context is not displayed")
		} else {
			t.Errorf("when the page displays an error: %s", res)
		}
	}
}

func TestRenderJson(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortPageRender, "test-render-json"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != `{"OK":true}` {
		if res == "" {
			t.Error("json is not displayed")
		} else {
			t.Errorf("when the json displays an error: %s", res)
		}
	}
}

func TestStrsliceFilter(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortPageRender, "test-strslice-filter"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "1, 2, 3|" {
		if res == "" {
			t.Error("slice is not displayed")
		} else {
			t.Errorf("when the slice displays an error: %s", res)
		}
	}
}

func TestCustomFilters(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortPageRender, "test-custom-filters"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "value1-filter1|value2-filter2" {
		if res == "" {
			t.Error("filters is not displayed")
		} else {
			t.Errorf("when the filters renders an error: %s", res)
		}
	}
}
