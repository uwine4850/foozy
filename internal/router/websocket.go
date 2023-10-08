package router

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websocket struct {
	upgrader      websocket.Upgrader
	conn          *websocket.Conn
	onMessage     func(messageType int, msgData []byte)
	onClientClose func()
}

func NewWebsocket(upgrader websocket.Upgrader) *Websocket {
	return &Websocket{upgrader: upgrader}
}

// Connect connecting web sockets to the client.
// fn func() is responsible for the event that occurs during the connection.
func (ws *Websocket) Connect(w http.ResponseWriter, r *http.Request, fn func()) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	ws.conn = conn
	fn()
	return nil
}

// Close closes the connection to the socket.
// Works well with the OnClientClose method.
func (ws *Websocket) Close() error {
	err := ws.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// OnClientClose event that will happen when the client closes the connection.
func (ws *Websocket) OnClientClose(fn func()) {
	ws.onClientClose = fn
}

// OnMessage event when a client sends a message.
func (ws *Websocket) OnMessage(fn func(messageType int, msgData []byte)) {
	ws.onMessage = fn
}

// SendMessage sending a message to the client.
func (ws *Websocket) SendMessage(messageType int, msg []byte) error {
	err := ws.conn.WriteMessage(messageType, msg)
	if err != nil {
		return err
	}
	return nil
}

// ReceiveMessages starts an infinite loop that listens for new messages from clients.
func (ws *Websocket) ReceiveMessages() error {
	for {
		messageType, msgData, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				ws.onClientClose()
				break
			}
			return err
		}
		ws.onMessage(messageType, msgData)
	}
	return nil
}
