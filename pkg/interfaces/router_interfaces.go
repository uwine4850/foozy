package interfaces

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

type IManager interface {
	itypeopr.INewInstance
	Render() IRender
	OneTimeData() IManagerOneTimeData
	Key() IKey
	Database() IDatabasePool
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

type IManagerOneTimeData interface {
	itypeopr.INewInstance
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
	DelUserContext(key string)
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
}

type IDatabasePool interface {
	ConnectionPool(name string) (IReadDatabase, error)
	AddConnection(name string, rd IReadDatabase) error
	Lock()
}

type IWebsocket interface {
	Close() error
	OnClientClose(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
	OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
	OnConnect(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
	SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
	ReceiveMessages(w http.ResponseWriter, r *http.Request) error
}

type IForm interface {
	Parse() error
	GetMultipartForm() *multipart.Form
	GetApplicationForm() url.Values
	Value(key string) string
	File(key string) (multipart.File, *multipart.FileHeader, error)
	Files(key string) ([]*multipart.FileHeader, bool)
}
