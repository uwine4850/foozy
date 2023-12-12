package interfaces

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type IManager interface {
	IManagerRender
	IManagerWebsocket
	IManagerData
}

type IManagerData interface {
	SetContext(data map[string]interface{})
	GetSlugParams(key string) (string, bool)
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
	DelUserContext(key string)
}

type IManagerRender interface {
	SetTemplateEngine(engine ITemplateEngine)
	RenderTemplate(w http.ResponseWriter, r *http.Request) error
	SetTemplatePath(templatePath string)
	SetSlugParams(params map[string]string)
	RenderJson(data interface{}, w http.ResponseWriter) error
}

type IManagerWebsocket interface {
	CurrentWebsocket() IWebsocket
	SetWebsocket(websocket IWebsocket)
}

type IWebsocket interface {
	Close() error
	OnClientClose(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
	OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
	OnConnect(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
	SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
	ReceiveMessages(w http.ResponseWriter, r *http.Request) error
}
