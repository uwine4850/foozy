package database

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/uwine4850/foozy/pkg/interfaces"
)

// DatabasePool stores a pool of database connections.
// To work properly, you need to connect to the database using the [Open] method
// and pass the open connection to the [ConnectionPool] method. There can be several connections,
// but usually one connection pool is enough.
//
// After all the settings you need to call the [Lock] method to lock all the changes. After calling
// this method you can no longer add connections to the pool.
// Thus, the connection pool remains static and predictable from the beginning to the end of the application.
type DatabasePool struct {
	// Type map[string]interfaces.IReadDatabase.
	connnectionPool sync.Map
	locked          atomic.Bool
}

func NewDatabasePool() *DatabasePool {
	return &DatabasePool{}
}

// ConnectionPool get a named connection pool.
func (dp *DatabasePool) ConnectionPool(name string) (interfaces.DatabaseInteraction, error) {
	if !dp.locked.Load() {
		return nil, &ErrDatabasePoolNotLocked{}
	}
	if val, ok := dp.connnectionPool.Load(name); ok {
		return val.(interfaces.DatabaseInteraction), nil
	} else {
		return nil, &ErrConnectionNotExists{Name: name}
	}
}

// AddConnection adding a new open named database connection pool.
func (dp *DatabasePool) AddConnection(name string, dbInteraction interfaces.DatabaseInteraction) error {
	if dp.locked.Load() {
		return &ErrDatabasePoolIsLocked{}
	}
	if _, ok := dp.connnectionPool.Load(name); ok {
		return &ErrConnectionAlreadyExists{Name: name}
	} else {
		dp.connnectionPool.Store(name, dbInteraction)
		return nil
	}
}

// Lock blocks further changes to the pool.
func (dp *DatabasePool) Lock() {
	dp.locked.Store(true)
}

type ErrConnectionAlreadyExists struct {
	Name string
}

func (e ErrConnectionAlreadyExists) Error() string {
	return fmt.Sprintf("a connection named %s already exists", e.Name)
}

type ErrConnectionNotExists struct {
	Name string
}

func (e ErrConnectionNotExists) Error() string {
	return fmt.Sprintf("a connection named %s not exists", e.Name)
}

type ErrDatabasePoolIsLocked struct{}

func (e ErrDatabasePoolIsLocked) Error() string {
	return "database pool is locked"
}

type ErrDatabasePoolNotLocked struct{}

func (e ErrDatabasePoolNotLocked) Error() string {
	return "database pool is not locked"
}
