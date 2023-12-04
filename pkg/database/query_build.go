package database

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"reflect"
	"strings"
)

var operators = []string{"!=", "=", "<", ">", "<=", ">=", "IN"}

type selectValues struct {
	Cols      string
	TableName string
}

type insertValues struct {
	TableName string
	Cols      string
	Values    string
	Args      []interface{}
}

type deleteValues struct {
	TableName string
}

type updateValues struct {
	TableName string
	StrVal    string
	Args      []interface{}
}

// QueryBuild builder of a database query.
// With this builder it is possible to create more flexible queries.
type QueryBuild struct {
	selectVal      selectValues
	insertVal      insertValues
	deleteVal      deleteValues
	updateVal      updateValues
	primaryCommand string
	queryStr       string
	queryArgs      []interface{}
	keyForAsyncQ   string
	syncQ          interfaces.ISyncQueries
	asyncQ         interfaces.IAsyncQueries
}

// SetSyncQ sets the synchronous query interface to communicate with the database.
func (qb *QueryBuild) SetSyncQ(sq interfaces.ISyncQueries) {
	qb.primaryCommand = ""
	qb.syncQ = sq
}

// SetAsyncQ sets the asynchronous query interface to communicate with the database.
func (qb *QueryBuild) SetAsyncQ(aq interfaces.IAsyncQueries) {
	qb.primaryCommand = ""
	qb.asyncQ = aq
}

// SetKeyForAsyncQ sets the key for asynchronous requests.
func (qb *QueryBuild) SetKeyForAsyncQ(key string) {
	qb.keyForAsyncQ = key
}

// Select database query.
func (qb *QueryBuild) Select(cols string, tableName string) interfaces.IUserQueryBuild {
	if qb.primaryCommand != "" {
		panic(fmt.Sprintf("you cannot use an %s command while already using a %s command.", "SELECT",
			qb.primaryCommand))
		return qb
	}
	qb.primaryCommand = "SELECT"
	qb.selectVal.Cols = cols
	qb.selectVal.TableName = tableName
	return qb
}

// Insert database query.
func (qb *QueryBuild) Insert(tableName string, params map[string]interface{}) interfaces.IUserQueryBuild {
	if qb.primaryCommand != "" {
		panic(fmt.Sprintf("you cannot use an %s command while already using a %s command.", "INSERT",
			qb.primaryCommand))
		return qb
	}
	qb.primaryCommand = "INSERT"
	keys, args := dbutils.ParseParams(params)
	qb.insertVal.TableName = tableName
	qb.insertVal.Cols = strings.Join(keys, ", ")
	qb.insertVal.Values = dbutils.RepeatValues(len(args), ",")
	qb.insertVal.Args = args
	return qb
}

// Delete database query.
// Deletes all table data if the Where method is not used.
func (qb *QueryBuild) Delete(tableName string) interfaces.IUserQueryBuild {
	if qb.primaryCommand != "" {
		panic(fmt.Sprintf("you cannot use an %s command while already using a %s command.", "DELETE",
			qb.primaryCommand))
		return qb
	}
	qb.primaryCommand = "DELETE"
	qb.deleteVal.TableName = tableName
	return qb
}

// Update database query.
func (qb *QueryBuild) Update(tableName string, params map[string]interface{}) interfaces.IUserQueryBuild {
	if qb.primaryCommand != "" {
		panic(fmt.Sprintf("you cannot use an %s command while already using a %s command.", "UPDATE",
			qb.primaryCommand))
		return qb
	}
	qb.primaryCommand = "UPDATE"
	var strVal string
	var args []interface{}
	i := 0
	for col, value := range params {
		args = append(args, value)
		if i == len(params)-1 {
			strVal += fmt.Sprintf("%s = ?", col)
		} else {
			strVal += fmt.Sprintf("%s = ?, ", col)
		}
		i++
	}
	qb.updateVal.TableName = tableName
	qb.updateVal.StrVal = strVal
	qb.updateVal.Args = args
	return qb
}

// Where adds a condition to the query.
// Example of using Where("id", "=", 55).
// For proper performance it is important to select the key character separately, the rest can be a single string.
// All characters to be selected are in the variable operators.
func (qb *QueryBuild) Where(args ...any) interfaces.IUserQueryBuild {
	var isNextValue bool
	var isNextIN bool
	var values []interface{}
	where := "WHERE "
	for i := 0; i < len(args); i++ {
		if isNextValue {
			isNextValue = false
			if isNextIN {
				isNextIN = false
				if reflect.TypeOf(args[i]).Kind() == reflect.Slice {
					slice := args[i].([]interface{})
					strVal := "(" + dbutils.RepeatValues(len(slice), ",") + ")"
					values = append(values, slice...)
					where += strVal + " "
				} else {
					panic("the IN operator must always be followed by a value of type []interface{}")
				}
				continue
			}
			where += "? "
			values = append(values, args[i])
			continue
		}
		if reflect.TypeOf(args[i]).Kind() == reflect.String && utils.SliceContains(operators, reflect.ValueOf(args[i]).String()) {
			isNextValue = true
			if reflect.ValueOf(args[i]).String() == "IN" {
				isNextIN = true
			}
			where += reflect.ValueOf(args[i]).String() + " "
			continue
		}
		if reflect.TypeOf(args[i]).Kind() == reflect.String {
			where += reflect.ValueOf(args[i]).String() + " "
		}
	}
	qb.queryStr = where + " "
	qb.queryArgs = values
	return qb
}

// Count outputs the number of records.
// Combined with Select only.
func (qb *QueryBuild) Count() interfaces.IUserQueryBuild {
	if qb.primaryCommand != "SELECT" {
		panic("the COUNT command can only be used with the SELECT command")
	}
	qb.selectVal.Cols = "COUNT(" + qb.selectVal.Cols + ")"
	return qb
}

// buildPrimaryCommand build a query to the database.
func (qb *QueryBuild) buildPrimaryCommand() (string, []interface{}) {
	var queryStr string
	var args []interface{}
	switch qb.primaryCommand {
	case "SELECT":
		queryStr += "SELECT " + qb.selectVal.Cols + " FROM " + qb.selectVal.TableName + " "
	case "INSERT":
		queryStr += fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ( %s ) ", qb.insertVal.TableName, qb.insertVal.Cols, qb.insertVal.Values)
		args = qb.insertVal.Args
	case "DELETE":
		queryStr += fmt.Sprintf("DELETE FROM %s ", qb.deleteVal.TableName)
	case "UPDATE":
		queryStr += fmt.Sprintf("UPDATE %s SET %s ", qb.updateVal.TableName, qb.updateVal.StrVal)
		args = qb.updateVal.Args
	}
	return queryStr, args
}

// Ex sending a query to the database and receiving the result (if the query is synchronous).
func (qb *QueryBuild) Ex() ([]map[string]interface{}, error) {
	queryStr, args := qb.buildPrimaryCommand()
	queryStr += qb.queryStr
	args = append(args, qb.queryArgs...)
	if qb.asyncQ != nil {
		qb.asyncQ.AsyncQuery(qb.keyForAsyncQ, queryStr, args...)
		return nil, nil
	} else {
		return qb.syncQ.Query(queryStr, args...)
	}
}
