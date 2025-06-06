package formtest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
	"github.com/uwine4850/foozy/tests/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests/common/tutils"
)

type Fill struct {
	NilField []string `form:"isNil"`
	Str      string
	Field1   []string        `form:"f1"`
	File     []form.FormFile `form:"file" empty:""`
}

func removeAllFilesInDirectory(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		filePath := dirPath + "/" + fileInfo.Name()

		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	testinitcnf.InitCnf()
	err := removeAllFilesInDirectory("./saved_files")
	if err != nil {
		panic(err)
	}

	mddl := middlewares.NewMiddleware()
	mddl.AsyncMddl(builtin_mddl.GenerateAndSetCsrf(1800, nil))
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newRouter := router.NewRouter(manager.NewManager(render))
	newRouter.SetMiddleware(mddl)
	newRouter.Post("/application-form", applicationForm)
	newRouter.Post("/multipart-form", multipartForm)
	newRouter.Post("/save-file", saveFile)
	serv := server.NewServer(tconf.PortForm, newRouter, nil)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortForm, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
}

func saveFile(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	_, header, err := newForm.File("file")
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	var path string
	err = form.SaveFile(w, header, "./saved_files", &path, manager)
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	if !fpath.PathExist(path) {
		return func() { w.Write([]byte("File not found.")) }
	}
	return func() {}
}

func multipartForm(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	_, header, err := newForm.File("file")
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	if header.Filename != "x.png" {
		return func() { w.Write([]byte("The file to be sent was not found.")) }
	}
	if newForm.Value("f1") != "v1" {
		return func() { w.Write([]byte("The field f1 to be sent was not found.")) }
	}
	return func() {}
}

func applicationForm(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	value1 := newForm.Value("f1")
	value2 := newForm.Value("f2")
	if value1 != "v1" || value2 != "v2" {
		return func() { w.Write([]byte("These forms do not match the forms sent.")) }
	}
	return func() {}
}

func TestApplicationForm(t *testing.T) {
	resp, err := form.SendApplicationForm(tutils.MakeUrl(tconf.PortForm, "application-form"), map[string][]string{"f1": {"v1"}, "f2": {"v2"}})
	if err != nil {
		t.Error(err)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
}

func TestMultipartForm(t *testing.T) {
	multipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortForm, "multipart-form"), map[string][]string{"f1": {"v1"}}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}

func TestSaveFile(t *testing.T) {
	sendMultipartForm, err := form.SendMultipartForm(tutils.MakeUrl(tconf.PortForm, "save-file"), map[string][]string{}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	defer sendMultipartForm.Body.Close()
	responseBody, err := io.ReadAll(sendMultipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Error(string(responseBody))
	}
}
