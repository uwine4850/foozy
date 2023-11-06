package interfaces

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type IManager interface {
	SetTemplateEngine(engine ITemplateEngine)
	RenderTemplate(w http.ResponseWriter, r *http.Request) error
	SetTemplatePath(templatePath string)
	SetContext(data map[string]interface{})
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
	GetWebSocket() IWebsocket
	SetWebsocket(websocket IWebsocket)
	RenderJson(data interface{}, w http.ResponseWriter) error
	DelUserContext(key string)
}

type IRouter interface {
	Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager) func())
	Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager) func())
	Ws(pattern string, ws IWebsocket, fn func(w http.ResponseWriter, r *http.Request, manager IManager) func())
	GetMux() *http.ServeMux
	SetTemplateEngine(engine ITemplateEngine)
	SetMiddleware(middleware IMiddleware)
}

type IWebsocket interface {
	Connect(w http.ResponseWriter, r *http.Request, fn func()) error
	Close() error
	OnClientClose(fn func(conn *websocket.Conn))
	OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
	OnConnect(fn func(conn *websocket.Conn))
	SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
	ReceiveMessages(w http.ResponseWriter, r *http.Request) error
}
