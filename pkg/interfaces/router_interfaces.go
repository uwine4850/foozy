package interfaces

import (
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
	Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	Ws(pattern string, ws IWebsocket, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	GetMux() *http.ServeMux
	SetTemplateEngine(engine ITemplateEngine)
	SetMiddleware(middleware IMiddleware)
}

type IWebsocket interface {
	Connect(w http.ResponseWriter, r *http.Request, fn func()) error
	Close() error
	OnClientClose(fn func())
	OnMessage(fn func(messageType int, msgData []byte))
	SendMessage(messageType int, msg []byte) error
	ReceiveMessages() error
}
