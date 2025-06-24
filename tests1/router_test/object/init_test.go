package object_test_1

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
)

type DatabaseTable struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Ok   bool   `db:"ok"`
}

type DTOMessage struct {
	rest.ImplementDTOMessage
	TypDTOMessage rest.TypeId `dto:"-typeid"`
	Id            int         `dto:"Id"`
	Name          string      `dto:"Name"`
	Ok            bool        `dto:"Ok"`
}

var newDTO = rest.NewDTO()
var messages = map[string][]irest.IMessage{
	"tee.ts": {
		DTOMessage{},
	},
}
var allowMessages = []rest.AllowMessage{
	{
		Package: "object_test_1",
		Name:    "DTOMessage",
	},
}

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()

	newDTO.Messages(messages)
	newDTO.AllowedMessages(allowMessages)

	sqlmock.New()
	newRender, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newManager := manager.NewManager(
		manager.NewOneTimeData(),
		newRender,
		manager.NewDatabasePool(),
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/test-object-view/:id", objectView())
	newRouter.Register(router.MethodGET, "/test-all-view", objectAllView())
	newRouter.Register(router.MethodGET, "/test-all-view/:id", objectAllSlugView())
	newRouter.Register(router.MethodGET, "/test-multiple-view/:o1/:o2", objectMultipleView())
	newRouter.Register(router.MethodGET, "/test-json-object-view/:id", jsonObjectView())
	newRouter.Register(router.MethodGET, "/test-json-all-view", objectJsonAllView())
	newRouter.Register(router.MethodGET, "/test-slug-json-all-view/:id", objectSlugJsonAllView())
	newRouter.Register(router.MethodGET, "/test-json-multiple-view/:o1/:o2", objectJsonMultipleView())
	newServer := server.NewServer(tutils.PortObject, newRouter, nil)
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
