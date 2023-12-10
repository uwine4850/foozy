package router

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
	MsgType int
	Msg     []byte
	Conn    *websocket.Conn
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websocket struct {
	upgrader      websocket.Upgrader
	conn          *websocket.Conn
	onMessage     func(messageType int, msgData []byte, conn *websocket.Conn)
	onClientClose func(conn *websocket.Conn)
	onConnect     func(conn *websocket.Conn)
	broadcast     chan Message
}

func NewWebsocket(upgrader websocket.Upgrader) *Websocket {
	return &Websocket{upgrader: upgrader, broadcast: make(chan Message)}
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
func (ws *Websocket) OnClientClose(fn func(conn *websocket.Conn)) {
	ws.onClientClose = fn
}

// OnMessage event when a client sends a message.
func (ws *Websocket) OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn)) {
	ws.onMessage = fn
}

func (ws *Websocket) OnConnect(fn func(conn *websocket.Conn)) {
	ws.onConnect = fn
}

// SendMessage sending a message to the client.
func (ws *Websocket) SendMessage(messageType int, msg []byte, conn *websocket.Conn) error {
	err := conn.WriteMessage(messageType, msg)
	if err != nil {
		return err
	}
	return nil
}

// ReceiveMessages starts an infinite loop that listens for new messages from clients.
func (ws *Websocket) ReceiveMessages(w http.ResponseWriter, r *http.Request) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	if ws.onConnect != nil {
		ws.onConnect(conn)
	}
	go ws.receiveMessages()
	for {
		messageType, msgData, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				if ws.onClientClose != nil {
					ws.onClientClose(conn)
				}
				break
			}
			return err
		}
		msg := Message{
			MsgType: messageType,
			Msg:     msgData,
			Conn:    conn,
		}
		ws.broadcast <- msg
	}
	return nil
}

// receiveMessages processing of received messages.
func (ws *Websocket) receiveMessages() {
	if ws.onMessage == nil {
		panic("OnMessage handler was not found")
	}
	for {
		msg := <-ws.broadcast
		ws.onMessage(msg.MsgType, msg.Msg, msg.Conn)
	}
}

func WsSendTextMessage(msg string, url string) (*http.Response, error) {
	dial, response, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return response, err
	}
	err = dial.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return response, err
	}
	err = dial.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return response, err
	}
	err = dial.Close()
	if err != nil {
		return response, err
	}
	return response, nil
}
