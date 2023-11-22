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
	QB() IQueryBuild
	SetDB(db *sql.DB)
	Query(query string, args ...any) ([]map[string]interface{}, error)
	Insert(tableName string, params map[string]interface{}) ([]map[string]interface{}, error)
	Select(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
	Update(tableName string, params []dbutils.DbEquals, where dbutils.WHOutput) ([]map[string]interface{}, error)
	Delete(tableName string, where dbutils.WHOutput) ([]map[string]interface{}, error)
	Count(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
}

type IAsyncQueries interface {
	QB(key string) IQueryBuild
	SetSyncQueries(queries ISyncQueries)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	AsyncQuery(key string, query string, args ...any)
	AsyncSelect(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
	AsyncInsert(key string, tableName string, params map[string]interface{})
	AsyncUpdate(key string, tableName string, params []dbutils.DbEquals, where dbutils.WHOutput)
	AsyncDelete(key string, tableName string, where dbutils.WHOutput)
	AsyncCount(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
}

type IQueryBuild interface {
	SetSyncQ(sq ISyncQueries)
	SetAsyncQ(aq IAsyncQueries)
	SetKeyForAsyncQ(key string)
	Select(cols string, tableName string) IQueryBuild
	Insert(tableName string, params map[string]interface{}) IQueryBuild
	Delete(tableName string) IQueryBuild
	Update(tableName string, params map[string]interface{}) IQueryBuild
	Where(args ...any) IQueryBuild
	Count() IQueryBuild
	Ex() ([]map[string]interface{}, error)
}
