package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/interfaces"
	router2 "github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	fserer "github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/utils"
)

type Fill struct {
	NilField []string `form:"isNil"`
	Str      string
	Field1   []string        `form:"f1"`
	File     []form.FormFile `form:"file"`
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
	err := removeAllFilesInDirectory("./saved_files")
	if err != nil {
		panic(err)
	}

	mddl := middlewares.NewMiddleware()
	mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newRouter := router2.NewRouter(manager.NewManager(render))
	newRouter.EnableLog(false)
	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.SetMiddleware(mddl)
	newRouter.Post("/application-form", applicationForm)
	newRouter.Post("/multipart-form", multipartForm)
	newRouter.Post("/save-file", saveFile)
	newRouter.Post("/fill", fill)
	server := fserer.NewServer(":8020", newRouter)
	go func() {
		err = server.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}()
	exitCode := m.Run()
	err = server.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func fill(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	var f Fill
	err = form.FillStructFromForm(newForm, form.NewFillableFormStruct(&f), []string{"isNil"})
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	if f.NilField != nil {
		return func() { w.Write([]byte("The NilField must be nil.")) }
	}
	if f.Field1 == nil {
		return func() { w.Write([]byte("The Field1 field must be populated.")) }
	}
	if f.File == nil {
		return func() { w.Write([]byte("The File field must be populated.")) }
	}
	return func() {}
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
	err = form.SaveFile(w, header, "./saved_files", &path, manager.Config())
	if err != nil {
		return func() { w.Write([]byte(err.Error())) }
	}
	if !utils.PathExist(path) {
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
	resp, err := form.SendApplicationForm("http://localhost:8020/application-form", map[string]string{"f1": "v1", "f2": "v2"})
	if err != nil {
		t.Error(err)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Errorf(string(responseBody))
	}
	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
}

func TestMultipartForm(t *testing.T) {
	multipartForm, err := form.SendMultipartForm("http://localhost:8020/multipart-form", map[string]string{"f1": "v1"}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	defer multipartForm.Body.Close()
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Errorf(string(responseBody))
	}
}

func TestSaveFile(t *testing.T) {
	sendMultipartForm, err := form.SendMultipartForm("http://localhost:8020/save-file", map[string]string{}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	defer sendMultipartForm.Body.Close()
	responseBody, err := io.ReadAll(sendMultipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Errorf(string(responseBody))
	}
}

func TestFillStructFromForm(t *testing.T) {
	multipartForm, err := form.SendMultipartForm("http://localhost:8020/fill", map[string]string{"f1": "v1"}, map[string][]string{"file": {"x.png"}})
	if err != nil {
		t.Error(err)
	}
	responseBody, err := io.ReadAll(multipartForm.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "" {
		t.Errorf(string(responseBody))
	}
	err = multipartForm.Body.Close()
	if err != nil {
		panic(err)
	}
}
