## websocket
Implementation of the websocket protocol. Uses the [gorilla/websocket](https://github.com/gorilla/websocket) library.

Used as a regular handler, for example:
```golang
func Socket(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	socket := router.NewWebsocket(router.Upgrader)
	socket.OnConnect(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
        fmt.Println("Connect")
    })
	socket.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		fmt.Println("Close")
	})
	socket.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		if err := conn.WriteMessage(messageType, msgData) {
			fmt.Println("Send message error:", err)
		}
	})
    // Start receive messages
	if err := socket.ReceiveMessages(w, r); err != nil {
		fmt.Println("Receive messages error:", err)
	}
	return nil
}
```
__NOTE:__ to start the websocket, you __must__ run the method [ReceiveMessages](#websocketreceivemessages).

### Websocket object
The object that handles the Websocket connection. To start working with Websocket, you need to initialize this object.
```golang
socket := router.NewWebsocket(router.Upgrader)
```

#### Websocket.OnConnect
A handler that will be executed when the user connects to the Websocket.
```golang
socket.OnConnect(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
    fmt.Println("Connect")
})
```

#### Websocket.OnClientClose
A handler that will be executed when the user disconnects from Websocket.
```golang
socket.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
    fmt.Println("Close")
})
```

#### Websocket.OnMessage
A handler that will be executed when a message is received from any connected user.

* messageType — message type
* msgData — message data
* conn — connection that sent the message
```golang
socket.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
	if err := conn.WriteMessage(messageType, msgData) {
		fmt.Println("Send message error:", err)
	}
})
```

#### Websocket.ReceiveMessages
A method that listens to the websocket and receives messages from users. It is __required__ to start the websocket.
```golang
if err := socket.ReceiveMessages(w, r); err != nil {
	fmt.Println("Receive messages error:", err)
}
```

#### WsSendTextMessage
Sends a message to an open websocket. This function can be useful for testing or communicating with other websockets.

__IMPORTANT:__ This function creates a new connection each time, so it should not be used under high loads, as performance may be significantly reduced.
```golang
func WsSendTextMessage(msg string, url string, header http.Header) (*http.Response, error) {
	dial, response, err := websocket.DefaultDialer.Dial(url, header)
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
```