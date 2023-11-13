package dbutils

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// AsyncQueryData a structure that represents the result of executing an asynchronous database query.
type AsyncQueryData struct {
	Res   []map[string]interface{}
	Error error
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

// ParseString processing text values from a database.
func ParseString(value interface{}) string {
	if value == nil {
		return ""
	}
	_uint8 := value.([]uint8)
	return string(_uint8)
}

// ParseInt processing integer values from a database.
// Considered, the interface format may be []uint8.
func ParseInt(value interface{}) (int, error) {
	if value == nil {
		return -1, errors.New("value is nil")
	}
	_type := reflect.TypeOf(value).String()
	var v int64
	switch _type {
	case "[]uint8":
		_uint8 := value.([]uint8)
		parseInt, err := strconv.ParseInt(string(_uint8), 0, 64)
		if err != nil {
			return -1, err
		}
		v = parseInt
	case "int64":
		v = value.(int64)
	}
	return int(v), nil
}

// ParseDateTime parsing a date value from a database.
// The date format must be set in the layout form.
// For example, if the date is 2023-04-06 08:04:05, the layout will be 2006-01-02 15:04:05.
// The time of the template should not change, only the form can change.
func ParseDateTime(layout string, value interface{}) (time.Time, error) {
	strValue := ParseString(value)
	if strValue == "" {
		return time.Now(), errors.New("value is nil")
	}
	parse, err := time.Parse(layout, strValue)
	if err != nil {
		return parse, err
	}
	return parse, nil
}

// ParseFloat processing the float data type from the database.
// It is taken into account that the interface format can be []uint8.
func ParseFloat(value interface{}) (float64, error) {
	if value == nil {
		return -1, errors.New("value is nil")
	}
	_type := reflect.TypeOf(value).String()
	var v float64
	switch _type {
	case "[]uint8":
		_uint8 := value.([]uint8)
		float, err := strconv.ParseFloat(string(_uint8), 64)
		if err != nil {
			return float, err
		}
		v = float
	case "float32":
		_v := value.(float32)
		v = float64(_v)
	}
	return v, nil
}
