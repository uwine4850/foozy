package form_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
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
	newRouter.Register(router.MethodPOST, "/test-application-form", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		frm := form.NewForm(r)
		if err := frm.Parse(); err != nil {
			return err
		}
		v1 := frm.Value("v1")
		v2 := frm.Value("v2")
		if v1 != "v1" && v2 != "v2" {
			return errors.New("form values are not as expected")
		}
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodPOST, "/test-multipart-form", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		frm := form.NewForm(r)
		if err := frm.Parse(); err != nil {
			return err
		}
		_, header, err := frm.File("file")
		if err != nil {
			return err
		}
		if header.Filename != "x.png" {
			return errors.New("file names do not match")
		}
		if frm.Value("v1") != "v1" {
			return errors.New("form values are not as expected")
		}
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodPOST, "/test-savefile", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		newForm := form.NewForm(r)
		err := newForm.Parse()
		if err != nil {
			return err
		}
		_, header, err := newForm.File("file")
		if err != nil {
			return err
		}
		var path string
		err = form.SaveFile(header, "./saved_files", &path, manager)
		if err != nil {
			return err
		}
		if !fpath.PathExist(path) {
			return errors.New("file not found")
		}
		w.Write([]byte("OK"))
		return nil
	})
	newServer := server.NewServer(tutils.PortForm, newRouter, nil)
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

func TestApplicationForm(t *testing.T) {
	resp, err := form.SendApplicationForm(tutils.MakeUrl(tutils.PortForm, "test-application-form"), map[string][]string{
		"v1": {"v1"},
		"v2": {"v2"},
	})
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("application form error: %s", res)
	}
}

func TestMultipartForm(t *testing.T) {
	resp, err := form.SendMultipartForm(tutils.MakeUrl(tutils.PortForm, "test-multipart-form"),
		map[string][]string{"v1": {"v1"}}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(res)
	}
	if res != "OK" {
		t.Errorf("multipart form error: %s", res)
	}
}

func TestSaveFile(t *testing.T) {
	resp, err := form.SendMultipartForm(tutils.MakeUrl(tutils.PortForm, "test-savefile"),
		map[string][]string{}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(res)
	}
	if res != "OK" {
		t.Errorf("save fiel error: %s", res)
	}
}
