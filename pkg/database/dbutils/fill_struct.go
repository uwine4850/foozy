package dbutils

import (
	"bytes"
	"fmt"
	"github.com/uwine4850/foozy/pkg/ferrors"
	"reflect"
)

func anyToBytes(value interface{}) ([]byte, error) {
	var buf bytes.Buffer

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("%v", val.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(fmt.Sprintf("%v", val.Uint()))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("%v", val.Float()))
	case reflect.String:
		buf.WriteString(val.String())
	default:
		return nil, ferrors.ErrUnsupportedTypeConvert{Type: val.Kind().String()}
	}
	return buf.Bytes(), nil
}

// FillStructFromDb Fills the structure with data from the database.
// Each variable of the structure to be filled must have a "db" tag which is responsible for the name of the column in
// the database, e.g. `db: "name"`.
func FillStructFromDb(dbRes map[string]interface{}, fill interface{}) error {
	if reflect.TypeOf(fill).Kind() != reflect.Ptr {
		return ferrors.ErrParameterNotPointer{Param: "fill"}
	}
	if reflect.TypeOf(fill).Elem().Kind() != reflect.Struct {
		return ferrors.ErrParameterNotStruct{Param: "fill"}
	}

	t := reflect.TypeOf(fill).Elem()
	v := reflect.ValueOf(fill).Elem()
	var res []byte
	var dbResField interface{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("db")
		// Skip if the tag is not a form
		if tag == "" {
			continue
		}
		__dbResField, ok := dbRes[tag]
		if !ok {
			return ErrDbResFieldNotFound{Field: tag}
		}
		dbResField = __dbResField
		if dbResField != nil {
			if reflect.TypeOf(dbResField).Kind() == reflect.Slice {
				res, ok = dbResField.([]uint8)
				if !ok {
					return ferrors.ErrConvertType{Type1: reflect.TypeOf(dbResField).String(), Type2: "[]uint8"}
				}
			} else {
				toBytes, err := anyToBytes(dbResField)
				if err != nil {
					return err
				}
				res = toBytes
			}
			value.SetString(string(res))
		}
	}
	return nil
}
