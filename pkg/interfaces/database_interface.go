package interfaces

import (
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

type Database interface {
	DatabaseInteraction
	Open() error
	Close() error
}

type SyncAsyncQuery interface {
	// SyncQ gives access to [SyncQueries].
	// Unlike [NewAsyncQ], it does not create any new instances,
	// but simply provides access to one shared instance.
	SyncQ() SyncQ
	// NewAsyncQ NewAsyncQ creates a new instance of [AsyncQueries]
	// for asynchronous queries. The new instance is created using
	// the [INewInstance] interface.
	NewAsyncQ() (AsyncQ, error)
}

type DatabaseInteraction interface {
	SyncAsyncQuery
	// NewTransaction creates a new transaction instance.
	// Creating a new [Transaction] object for each transaction
	// is necessary for data security reasons.
	// That is, for each new transaction a new instance of [Transaction]
	// should be created, which works in its own scope. This behavior prevents any data leaks.
	NewTransaction() (DatabaseTransaction, error)
}

type DatabaseTransaction interface {
	SyncAsyncQuery
	NewAsyncQ() (AsyncQ, error)
	BeginTransaction() error
	CommitTransaction() error
	RollBackTransaction() error
}

// QueryExec an interface represents any object that can query a database.
type QueryExec interface {
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

type SyncQ interface {
	itypeopr.NewInstance
	QueryExec
	SetDB(db QueryExec)
}

type AsyncQ interface {
	itypeopr.NewInstance
	SetSyncQueries(queries SyncQ)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	Clear()
	Query(key string, query string, args ...any)
	Exec(key string, query string, args ...any)
}
