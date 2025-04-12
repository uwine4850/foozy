package interfaces

import (
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

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
	SetSyncQueries(queries ISyncQueries)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	Clear()
	Query(key string, query string, args ...any)
	Exec(key string, query string, args ...any)
}
