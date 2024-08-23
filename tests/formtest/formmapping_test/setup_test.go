package formmappingtest_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
)

func TestMain(m *testing.M) {
	mddl := middlewares.NewMiddleware()
	mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newRouter := router.NewRouter(manager.NewManager(render))
	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.SetMiddleware(mddl)
	newRouter.Post("/mp-default-struct", mpDefaultStruct)
	newRouter.Post("/mp-empty-string-0-err", mpEmptyString0Err)
	newRouter.Post("/mp-empty-string-1-err", mpEmptyString1Err)
	newRouter.Post("/mp-empty-file-err", mpEmptyFileErr)
	newRouter.Post("/mp-empty-value", mpEmptyValue)
	newRouter.Post("/fill", fill)
	newRouter.Post("/fill-reflect-value", fillReflectValue)
	newRouter.Post("/mp-typed-struct", mpTypedMapper)
	serv := server.NewServer(":8020", newRouter)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(":8020", 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}
