package interfaces

import (
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

type IDatabase interface {
	IReadDatabase
	IOpenCloseDatabase
}

type ISyncAsync interface {
	// SyncQ gives access to [ISyncQueries].
	// Unlike [NewAsyncQ], it does not create any new instances,
	// but simply provides access to one shared instance.
	SyncQ() ISyncQueries
	// NewAsyncQ NewAsyncQ creates a new instance of [IAsyncQueries]
	// for asynchronous queries. The new instance is created using
	// the [INewInstance] interface.
	NewAsyncQ() (IAsyncQueries, error)
}

type IReadDatabase interface {
	ISyncAsync
	// NewTransaction creates a new transaction instance.
	// Creating a new [ITransaction] object for each transaction
	// is necessary for data security reasons.
	// That is, for each new transaction a new instance of [ITransaction]
	// should be created, which works in its own scope. This behavior prevents any data leaks.
	NewTransaction() ITransaction
}

type ITransaction interface {
	ISyncAsync
	BeginTransaction() error
	CommitTransaction() error
	RollBackTransaction() error
}

type IOpenCloseDatabase interface {
	Open() error
	Close() error
}

// IDbQuery an interface represents any object that can query a database.
type IDbQuery interface {
	// Used to execute queries that return data.
	// For example, the SELECT command.
	Query(query string, args ...any) ([]map[string]interface{}, error)
	// Used to execute queries that do not return data.
	// For example, the INSERT command.
	//
	// Returns the following data:
	// Key "insertID" is the identifier of the inserted row using INSERT.
	// Key "rowsAffected" - Returns the number of rows affected by INSERT, UPDATE, DELETE.
	Exec(query string, args ...any) (map[string]interface{}, error)
}

type ISyncQueries interface {
	SetDB(db IDbQuery)
	Query(query string, args ...any) ([]map[string]interface{}, error)
	Exec(query string, args ...any) (map[string]interface{}, error)
}

type IAsyncQueries interface {
	itypeopr.INewInstance
	SetSyncQueries(queries ISyncQueries)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	Clear()
	Query(key string, query string, args ...any)
	Exec(key string, query string, args ...any)
}
