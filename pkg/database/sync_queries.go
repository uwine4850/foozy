package database

import (
	"database/sql"
	"fmt"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"strconv"
	"strings"
)

type SyncQueries struct {
	db *sql.DB
}

// Query sends a query to the database.
// The result of query execution is processed and converted to []map[string]interface{} format.
// The map key is the column names. The key values are the current column and string data in the interface{} format,
// which can be converted to the desired type.
func (q *SyncQueries) Query(query string, args ...any) ([]map[string]interface{}, error) {
	_query, err := q.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var rows []map[string]interface{}
	dbutils.ScanRows(_query, func(row map[string]interface{}) {
		rows = append(rows, row)
	})

	err = _query.Close()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (q *SyncQueries) SetDB(db *sql.DB) {
	q.db = db
}

func (q *SyncQueries) Select(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error) {
	var whereStr string
	if len(where.QueryArgs) > 0 {
		whereStr = " WHERE " + where.QueryStr
	}
	var limitStr string
	if limit > 0 {
		limitStr = "LIMIT " + strconv.Itoa(limit)
	}
	queryStr := fmt.Sprintf("SELECT %s FROM %s %s %s", strings.Join(rows, ", "), tableName, whereStr, limitStr)
	return q.Query(queryStr, where.QueryArgs...)
}

func (q *SyncQueries) Insert(tableName string, params map[string]interface{}) ([]map[string]interface{}, error) {
	keys, vals := dbutils.ParseParams(params)
	queryStr := fmt.Sprintf("INSERT INTO `%s` ( %s ) VALUES ( %s )",
		tableName, strings.Join(keys, ", "), dbutils.RepeatValues(len(vals), ","))
	return q.Query(queryStr, vals...)
}

func (q *SyncQueries) Delete(tableName string, where dbutils.WHOutput) ([]map[string]interface{}, error) {
	var whereStr string
	if len(where.QueryArgs) > 0 {
		whereStr = " WHERE " + where.QueryStr
	}
	queryStr := fmt.Sprintf("DELETE FROM %s %s", tableName, whereStr)
	return q.Query(queryStr, where.QueryArgs...)
}

func (q *SyncQueries) Update(tableName string, params []dbutils.DbEquals, where dbutils.WHOutput) ([]map[string]interface{}, error) {
	equalsStr, paramValues := dbutils.ParseEquals(params, ",")
	var whereStr string
	if len(where.QueryArgs) > 0 {
		whereStr = " WHERE " + where.QueryStr
	}
	queryStr := fmt.Sprintf("UPDATE `%s` SET %s %s ",
		tableName, equalsStr, whereStr)
	args := append(paramValues, where.QueryArgs...)
	return q.Query(queryStr, args...)
}

// Count returns the number of records under the condition.
// IMPORTANT: the result is in string format.
func (q *SyncQueries) Count(rows []string, tableName string, where dbutils.WHOutput, limit int) ([]map[string]interface{}, error) {
	var whereStr string
	if len(where.QueryArgs) > 0 {
		whereStr = " WHERE " + where.QueryStr
	}
	var limitStr string
	if limit > 0 {
		limitStr = "LIMIT " + strconv.Itoa(limit)
	}
	queryStr := fmt.Sprintf("SELECT COUNT(%s) FROM %s %s %s", strings.Join(rows, ", "), tableName, whereStr, limitStr)
	return q.Query(queryStr, where.QueryArgs...)
}
