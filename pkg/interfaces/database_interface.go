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
	// Key “id” is the identifier of the inserted row using INSERT.
	// Key “rows” - Returns the number of rows affected by INSERT, UPDATE, DELETE.
	Exec(query string, args ...any) (map[string]interface{}, error)
}

type ISyncQueries interface {
	SetDB(db IDbQuery)
	Query(query string, args ...any) ([]map[string]interface{}, error)
	Exec(query string, args ...any) (map[string]interface{}, error)
	Insert(tableName string, params map[string]any) (map[string]interface{}, error)
	Select(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
	Update(tableName string, params map[string]any, where dbutils.WHOutput) (map[string]interface{}, error)
	Delete(tableName string, where dbutils.WHOutput) (map[string]interface{}, error)
	Count(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
	Increment(fieldName string, tableName string, where dbutils.WHOutput) (map[string]interface{}, error)
	Exists(tableName string, where dbutils.WHOutput) (map[string]interface{}, error)
}

type IAsyncQueries interface {
	SetSyncQueries(queries ISyncQueries)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	Clear()
	Query(key string, query string, args ...any)
	Exec(key string, query string, args ...any)
	Select(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
	Insert(key string, tableName string, params map[string]any)
	Update(key string, tableName string, params map[string]any, where dbutils.WHOutput)
	Delete(key string, tableName string, where dbutils.WHOutput)
	Count(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
	Increment(key string, fieldName string, tableName string, where dbutils.WHOutput)
	Exists(key string, tableName string, where dbutils.WHOutput)
}
