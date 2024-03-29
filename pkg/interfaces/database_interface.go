package interfaces

import (
	"database/sql"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

type ISyncQueries interface {
	QB() IUserQueryBuild
	SetDB(db *sql.DB)
	Query(query string, args ...any) ([]map[string]interface{}, error)
	Insert(tableName string, params map[string]interface{}) ([]map[string]interface{}, error)
	Select(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
	Update(tableName string, params []dbutils.DbEquals, where dbutils.WHOutput) ([]map[string]interface{}, error)
	Delete(tableName string, where dbutils.WHOutput) ([]map[string]interface{}, error)
	Count(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error)
	Increment(fieldName string, tableName string, where dbutils.WHOutput) ([]map[string]interface{}, error)
}

type IAsyncQueries interface {
	QB(key string) IUserQueryBuild
	SetSyncQueries(queries ISyncQueries)
	Wait()
	LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
	AsyncQuery(key string, query string, args ...any)
	AsyncSelect(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
	AsyncInsert(key string, tableName string, params map[string]interface{})
	AsyncUpdate(key string, tableName string, params []dbutils.DbEquals, where dbutils.WHOutput)
	AsyncDelete(key string, tableName string, where dbutils.WHOutput)
	AsyncCount(key string, rows []string, tableName string, where dbutils.WHOutput, limit int)
	AsyncIncrement(key string, fieldName string, tableName string, where dbutils.WHOutput)
}

type IQueryBuild interface {
	IUserQueryBuild
	SetSyncQ(sq ISyncQueries)
	SetAsyncQ(aq IAsyncQueries)
	SetKeyForAsyncQ(key string)
}

type IUserQueryBuild interface {
	Select(cols string, tableName string) IUserQueryBuild
	Insert(tableName string, params map[string]interface{}) IUserQueryBuild
	Delete(tableName string) IUserQueryBuild
	Update(tableName string, params map[string]interface{}) IUserQueryBuild
	Increment(fieldName string, tableName string) IUserQueryBuild
	Where(args ...any) IUserQueryBuild
	Count() IUserQueryBuild
	Ex() ([]map[string]interface{}, error)
}
