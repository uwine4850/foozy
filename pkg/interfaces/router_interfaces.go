package interfaces

import (
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

type Manager interface {
	itypeopr.NewInstance
	Render() Render
	OneTimeData() ManagerOneTimeData
	Key() Key
	Database() DatabasePool
}

type Key interface {
	HashKey() string
	OldHashKey() string
	BlockKey() string
	OldBlockKey() string
	StaticKey() string
	Date() time.Time
	GenerateBytesKeys(length int)
	Generate32BytesKeys()
	Get32BytesKey() Key
}

type ManagerOneTimeData interface {
	itypeopr.NewInstance
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
	DelUserContext(key string)
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
}

type DatabasePool interface {
	ConnectionPool(name string) (DatabaseInteraction, error)
	AddConnection(name string, rd DatabaseInteraction) error
	Lock()
}
