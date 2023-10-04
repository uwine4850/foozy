package dbutils

import (
	"database/sql"
	"fmt"
)

// AsyncQueryData a structure that represents the result of executing an asynchronous database query.
type AsyncQueryData struct {
	Res   []map[string]interface{}
	Error string
}

// DbEquals structure used to represent the column and the value to which it should be equal.
type DbEquals struct {
	Name  string
	Value interface{}
}

// RepeatValues repeats the "?" sign several times.
func RepeatValues(count int, sep string) string {
	var val string
	for i := 0; i < count; i++ {
		if i == count-1 {
			val += "?"
		} else {
			val += "?" + sep + " "
		}
	}
	return val
}

// ScanRows scans the rows that the executed database query provides.
// According to the number of columns it creates a map of interfaces to be filled. Then fills the map with the value of
// the row and places it in the slice.
func ScanRows(rows *sql.Rows, fn func(row map[string]interface{})) {
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	dataColumns := make([]interface{}, len(columns))
	for i := range columns {
		dataColumns[i] = new(interface{})
	}
	for rows.Next() {
		row := make(map[string]interface{})
		err := rows.Scan(dataColumns...)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(columns); i++ {
			v := dataColumns[i].(*interface{})
			row[columns[i]] = *v
		}
		fn(row)
	}
}

func ParseParams(params map[string]interface{}) ([]string, []interface{}) {
	var values []interface{}
	var keys []string
	for key, val := range params {
		keys = append(keys, key)
		values = append(values, val)
	}
	return keys, values
}

// ParseEquals handles a slice of the DbEquals structure.
// Converts the data into the format "key = ?", and then creates from this a part of the string for sql query.
// The conjunction parameter is responsible for the delimiter between the data, if there is more than one.
func ParseEquals(equals []DbEquals, conjunction string) (string, []interface{}) {
	var w string
	var values []interface{}
	for i := 0; i < len(equals); i++ {
		values = append(values, equals[i].Value)
		if i == len(equals)-1 {
			w += fmt.Sprintf("%s = ?", equals[i].Name)
		} else {
			w += fmt.Sprintf("%s = ? %s ", equals[i].Name, conjunction)
		}
	}
	return w, values
}
