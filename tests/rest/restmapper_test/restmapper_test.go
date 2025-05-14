package restmappertest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/tests/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests/common/tutils"
)

var mng = manager.NewManager(nil)

var newRouter = router.NewRouter(mng)
var dto = rest.NewDTO()

func TestMain(m *testing.M) {
	testinitcnf.InitCnf()
	dto.AllowedMessages([]rest.AllowMessage{
		{
			Name:    "JsonData",
			Package: "restmappertest",
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
	serv := server.NewServer(tconf.PortRestMapper, newRouter, nil)
	go func() {
		err := serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortRestMapper, 5); err != nil {
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
	Id       int                 `json:"Id" dto:"Id"`
	Name     string              `json:"Name" dto:"Name"`
	Slice    []string            `json:"Slice" dto:"Slice"`
	IsOk     bool                `json:"IsOk" dto:"IsOk"`
	Map      map[string]int      `json:"Map" dto:"Map"`
	SliceMap []map[string]string `json:"SliceMap" dto:"SliceMap"`
	MapSlice map[string][]int    `json:"MapSlice" dto:"MapSlice"`
}

func TestRestJson(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortRestMapper, "json"))
	if err != nil {
		t.Error(err)
	}

	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &jsonMap); err != nil {
		t.Error(err)
	}

	var jsonData JsonData
	if err := mapper.JsonToDTOMessage(jsonMap, dto, &jsonData); err != nil {
		t.Error(err)
	}
	fmt.Println(jsonData)
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
