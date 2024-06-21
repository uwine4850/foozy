## Package websocket
The web socket interface is implemented using the __github.com/gorilla/websocket__ library.
In the package ``router`` there is a global variable ``Upgrader`` which is required 
websocket operations.

__OnConnect__
```
OnConnect(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
```
The function that is launched when connecting to the client.

__Close__
```
Close() error
```
Closing the connection.

__OnClientClose__
```
OnClientClose(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
```
The function that will be executed when the client closes the connection.

__OnMessage__
```
OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
```
When the socket receives the message, the function ``fn'' will be executed.

__SendMessage__
```
SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
```
Sending a message to the client.

__ReceiveMessages__
```
ReceiveMessages(w http.ResponseWriter, r *http.Request) error
```
A method that starts receiving messages. This method must be running.