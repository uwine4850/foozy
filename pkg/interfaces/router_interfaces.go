package interfaces

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type IManager interface {
	IManagerRender
	IManagerWebsocket
	IManagerConfig
	Render() IManagerRender
	WS() IManagerWebsocket
	OneTimeData() IManagerOneTimeData
	SetOneTimeData(manager IManagerOneTimeData)
	Config() IManagerConfig
}

type IKey interface {
	HashKey() string
	OldHashKey() string
	BlockKey() string
	OldBlockKey() string
	StaticKey() string
	Date() time.Time
	GenerateBytesKeys(length int)
}

type IManagerConfig interface {
	Debug(enable bool)
	IsDebug() bool
	ErrorLogging(enable bool)
	IsErrorLogging() bool
	ErrorLoggingFile(path string)
	GetErrorLoggingFile() string
	Generate32BytesKeys()
	Get32BytesKey() IKey
}

type IManagerOneTimeData interface {
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
	DelUserContext(key string)
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
}

type IManagerRender interface {
	SetContext(data map[string]interface{})
	SetTemplateEngine(engine ITemplateEngine)
	RenderTemplate(w http.ResponseWriter, r *http.Request) error
	SetTemplatePath(templatePath string)
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
