package interfaces

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

type IManager interface {
	Render() IRender
	SetRender(render IRender)
	OneTimeData() IManagerOneTimeData
	SetOneTimeData(manager IManagerOneTimeData)
	SetConfig(cnf IManagerConfig)
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
	Generate32BytesKeys()
	Get32BytesKey() IKey
}

type IManagerDebugConfig interface {
	Debug(enable bool)
	IsDebug() bool
	ErrorLogging(enable bool)
	IsErrorLogging() bool
	ErrorLoggingFile(path string)
	GetErrorLoggingFile() string
	SkipLoggingLevel(skip int)
	LoggingLevel() int
}

type IManagerConfig interface {
	DebugConfig() IManagerDebugConfig
	PrintLog(enable bool)
	IsPrintLog() bool
	Key() IKey
}

type IManagerOneTimeData interface {
	itypeopr.INewInstance
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
	DelUserContext(key string)
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
}

type IWebsocket interface {
	Close() error
	OnClientClose(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
	OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
	OnConnect(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
	SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
	ReceiveMessages(w http.ResponseWriter, r *http.Request) error
}
