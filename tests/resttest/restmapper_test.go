package resttest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/router/rest/restmapper"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/typeopr"
	initcnf "github.com/uwine4850/foozy/tests/init_cnf"
)

var mng = manager.NewManager(nil)

var newRouter = router.NewRouter(mng)
var dto = rest.NewDTO()

func TestMain(m *testing.M) {
	initcnf.InitCnf()
	dto.AllowedMessages([]rest.AllowMessage{
		{
			Name:    "JsonData",
			Package: "resttest",
		},
	})
	newRouter.Get("/json", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		jsonData := JsonData{
			Id:       1,
			Name:     "name",
			Slice:    []string{"a1", "a2"},
			IsOk:     true,
			Map:      map[string]int{"1": 1, "2": 2, "3": 3},
			SliceMap: []map[string]string{{"m1": "1", "m11": "11"}, {"m2": "2", "m22": "22"}},
			MapSlice: map[string][]int{"m": {1, 2, 3}, "mm": {11, 22, 33}},
		}
		return func() {
			router.SendJson(jsonData, w)
		}
	})
	serv := server.NewServer(":8070", newRouter, nil)
	go func() {
		err := serv.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(":8070", 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err := serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

type JsonData struct {
	rest.ImplementDTOMessage
	Id       int                 `json:"Id"`
	Name     string              `json:"Name"`
	Slice    []string            `json:"Slice"`
	IsOk     bool                `json:"IsOk"`
	Map      map[string]int      `json:"Map"`
	SliceMap []map[string]string `json:"SliceMap"`
	MapSlice map[string][]int    `json:"MapSlice"`
}

func TestRestJson(t *testing.T) {
	get, err := http.Get("http://localhost:8070/json")
	if err != nil {
		t.Error(err)
	}

	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	var jsonMap map[string]interface{}
	if err := restmapper.JsonStringToMap(string(body), &jsonMap); err != nil {
		t.Error(err)
	}
	var jsonData JsonData
	if err := restmapper.JsonToMessage(&jsonMap, dto, typeopr.Ptr{}.New(&jsonData)); err != nil {
		t.Error(err)
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
