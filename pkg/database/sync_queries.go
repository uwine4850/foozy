package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type SyncQueries struct {
	db interfaces.IDbQuery
	qb interfaces.IQueryBuild
}

func NewSyncQueries(qb interfaces.IQueryBuild) *SyncQueries {
	return &SyncQueries{qb: qb}
}

func (q *SyncQueries) QB() interfaces.IUserQueryBuild {
	q.qb.SetSyncQ(q)
	return q.qb
}

// Query sends a query to the database.
// The result of query execution is processed and converted to []map[string]interface{} format.
// The map key is the column names. The key values are the current column and string data in the interface{} format,
// which can be converted to the desired type.
func (q *SyncQueries) Query(query string, args ...any) ([]map[string]interface{}, error) {
	return q.db.Query(query, args...)
}

func (q *SyncQueries) Exec(query string, args ...any) (map[string]interface{}, error) {
	return q.db.Exec(query, args...)
}

func (q *SyncQueries) SetDB(db interfaces.IDbQuery) {
	q.db = db
}

func (q *SyncQueries) Select(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error) {
	if len(rows) == 0 {
		return nil, errors.New("the rows argument cannot be empty")
	}
	var whereStr string
	if where.QueryStr != "" {
		whereStr = " WHERE " + where.QueryStr
	}
	var limitStr string
	if limit > 0 {
		limitStr = "LIMIT " + strconv.Itoa(limit)
	}
	queryStr := fmt.Sprintf("SELECT %s FROM %s %s %s", strings.Join(rows, ", "), tableName, whereStr, limitStr)
	return q.Query(queryStr, where.QueryArgs...)
}

func (q *SyncQueries) Insert(tableName string, params map[string]any) (map[string]interface{}, error) {
	keys, vals := dbutils.ParseParams(params)
	queryStr := fmt.Sprintf("INSERT INTO `%s` ( %s ) VALUES ( %s )",
		tableName, strings.Join(keys, ", "), dbutils.RepeatValues(len(vals), ","))
	return q.Exec(queryStr, vals...)
}

func (q *SyncQueries) Delete(tableName string, where dbutils.WHOutput) (map[string]interface{}, error) {
	var whereStr string
	if where.QueryStr != "" {
		whereStr = " WHERE " + where.QueryStr
	}
	queryStr := fmt.Sprintf("DELETE FROM %s %s", tableName, whereStr)
	return q.Exec(queryStr, where.QueryArgs...)
}

func (q *SyncQueries) Update(tableName string, params map[string]any, where dbutils.WHOutput) (map[string]interface{}, error) {
	equalsStr, paramValues := dbutils.ParseMapAsEquals(&params)
	var whereStr string
	if where.QueryStr != "" {
		whereStr = " WHERE " + where.QueryStr
	}
	queryStr := fmt.Sprintf("UPDATE `%s` SET %s %s ",
		tableName, equalsStr, whereStr)
	args := append(paramValues, where.QueryArgs...)
	return q.Exec(queryStr, args...)
}

// Count returns the number of records under the condition.
func (q *SyncQueries) Count(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error) {
	var whereStr string
	if where.QueryStr != "" {
		whereStr = " WHERE " + where.QueryStr
	}
	var limitStr string
	if limit > 0 {
		limitStr = "LIMIT " + strconv.Itoa(limit)
	}
	queryStr := fmt.Sprintf("SELECT COUNT(%s) FROM %s %s %s", strings.Join(rows, ", "), tableName, whereStr, limitStr)
	res, err := q.Query(queryStr, where.QueryArgs...)
	if err != nil {
		return nil, err
	}
	if err := dbutils.DatabaseResultNotEmpty(res); err != nil {
		return nil, err
	}
	return []map[string]interface{}{{"count": res[0][fmt.Sprintf("COUNT(%s)", strings.Join(rows, ", "))]}}, nil
}

// Increment does an increment of a field of type INT.
func (q *SyncQueries) Increment(fieldName string, tableName string, where dbutils.WHOutput) (map[string]interface{}, error) {
	var whereStr string
	if where.QueryStr != "" {
		whereStr = " WHERE " + where.QueryStr
	}
	queryStr := fmt.Sprintf("UPDATE `%s` SET `%s`= `%s`+ 1 %s ", tableName, fieldName, fieldName, whereStr)
	return q.Exec(queryStr, where.QueryArgs...)
}

func (q *SyncQueries) Exists(tableName string, where dbutils.WHOutput) (map[string]interface{}, error) {
	var whereStr string
	if where.QueryStr != "" {
		whereStr = " WHERE " + where.QueryStr
	}
	queryStr := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s %s) AS is_exists", tableName, whereStr)
	res, err := q.Query(queryStr, where.QueryArgs...)
	if err != nil {
		return nil, err
	}
	exists := res[0]["is_exists"].(int64) != 0
	return map[string]interface{}{"exists": exists}, nil
}
