package interfaces

import (
	"database/sql"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

type IDatabase interface {
	Connect() error
	Close() error
	SetSyncQueries(q ISyncQueries)
	SetAsyncQueries(q IAsyncQueries)
	SyncQ() ISyncQueries
	AsyncQ() IAsyncQueries
	DatabaseName() string
}

type ISyncQueries interface {
	SetDB(db *sql.DB)
	Query(query string, args ...any) ([]map[string]interface{}, error)
	Insert(tableName string, params map[string]interface{}) ([]map[string]interface{}, error)
	Select(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
	Update(tableName string, params []dbutils.DbEquals, where dbutils.WHOutput) ([]map[string]interface{}, error)
	Delete(tableName string, where dbutils.WHOutput) ([]map[string]interface{}, error)
	Count(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
}

type IAsyncQueries interface {
	SetSyncQueries(queries ISyncQueries)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	AsyncSelect(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
	AsyncInsert(key string, tableName string, params map[string]interface{})
	AsyncUpdate(key string, tableName string, params []dbutils.DbEquals, where dbutils.WHOutput)
	AsyncDelete(key string, tableName string, where dbutils.WHOutput)
	AsyncCount(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
}
