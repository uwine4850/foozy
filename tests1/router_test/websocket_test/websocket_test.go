package websocket_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
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
		manager.NewDatabasePool(),
	)
	newMiddlewares := middlewares.NewMiddlewares()
	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/socket", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		socket := router.NewWebsocket(router.Upgrader)
		socket.OnConnect(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {})
		socket.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {})
		socket.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
			if err := socket.SendMessage(messageType, msgData, conn); err != nil {
				panic(err)
			}
		})
		if err := socket.ReceiveMessages(w, r); err != nil {
			panic(err)
		}
		return nil
	})
	newServer := server.NewServer(tutils.PortSocket, newRouter, nil)
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

func TestSocket(t *testing.T) {
	connect, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost%s/socket", tutils.PortSocket), nil)
	if err != nil {
		t.Error(err)
	}
	sendData := map[string]string{"Message": "TEST"}
	if err := connect.WriteJSON(sendData); err != nil {
		t.Error(err)
	}
	_, message, err := connect.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	var outData map[string]string
	if err := json.Unmarshal(message, &outData); err != nil {
		t.Error(err)
	}
	if !maps.Equal(sendData, outData) {
		t.Error("sent and received data do not match")
	}
}
