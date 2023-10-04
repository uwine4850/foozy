package dbutils

import (
	"database/sql"
	"fmt"
)

type AsyncQueryData struct {
	Res   []map[string]interface{}
	Error string
}

type DbEquals struct {
	Name  string
	Value interface{}
}

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
			//rv := string((*v).([]uint8))
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

func ParseEquals(where []DbEquals, conjunction string) (string, []interface{}) {
	var w string
	var values []interface{}
	for i := 0; i < len(where); i++ {
		values = append(values, where[i].Value)
		if i == len(where)-1 {
			w += fmt.Sprintf("%s = ?", where[i].Name)
		} else {
			w += fmt.Sprintf("%s = ? %s ", where[i].Name, conjunction)
		}
	}
	return w, values
}
