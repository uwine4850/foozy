package dbmapper

import (
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

// FillReflectValueFromDb fills structure of type reflect.Value with data from a database query.
// For proper operation, you need to add a "name" tag with the column name for each field that represents a column from the database.
// For example: `name:"col_name"`.
// You can also use the "empty" tag to indicate the action when the field is empty. If this tag is empty, nothing will happen. Other meanings:
// Any text - sets the field values ​​to this text.
// -error - display an error if the field is empty.
func FillReflectValueFromDb(dbRes map[string]interface{}, fill *reflect.Value) error {
	t := fill.Type()
	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(namelib.TAGS.DB_MAPPER_NAME)
		if exist, err := validateDbTag(tagName, &dbRes); err != nil {
			return err
		} else {
			if !exist {
				continue
			}
		}
		val, err := dbValueConversionToByte(dbRes[tagName])
		if err != nil {
			return err
		}
		if typeopr.IsEmpty(val) {
			if err := emptyOperations(&val, t.Field(i).Tag.Get(namelib.TAGS.DB_MAPPER_EMPTY), tagName); err != nil {
				return err
			}
		}
		fill.FieldByName(t.Field(i).Name).SetString(string(val))
	}
	return nil
}

// FillMapFromDb fills the map that is passed by reference with values from the database.
func FillMapFromDb(dbRes map[string]interface{}, fill *map[string]string) error {
	if *fill == nil {
		panic("The \"fill\" map must not be of type nil.")
	}
	for key, value := range dbRes {
		if value == nil {
			(*fill)[key] = ""
			continue
		}

		val, err := dbValueConversionToByte(value)
		if err != nil {
			return err
		}
		(*fill)[key] = string(val)
	}
	return nil
}

// FillStructFromDb Fills the structure with data from the database.
// For proper operation, you need to add a "name" tag with the column name for each field that represents a column from the database.
// For example: `name:"col_name"`.
// You can also use the "empty" tag to indicate the action when the field is empty. If this tag is empty, nothing will happen. Other meanings:
// Any text - sets the field values ​​to this text.
// -error - display an error if the field is empty.
func FillStructFromDb(dbRes map[string]interface{}, fillPtr typeopr.IPtr) error {
	fill := fillPtr.Ptr()
	if !typeopr.PtrIsStruct(fill) {
		return typeopr.ErrParameterNotStruct{Param: "fill"}
	}

	t := reflect.TypeOf(fill).Elem()
	v := reflect.ValueOf(fill).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get(namelib.TAGS.DB_MAPPER_NAME)
		if exist, err := validateDbTag(tag, &dbRes); err != nil {
			return err
		} else {
			if !exist {
				continue
			}
		}
		dbResField := dbRes[tag]
		if dbResField != nil {
			val, err := dbValueConversionToByte(dbRes[tag])
			if err != nil {
				return err
			}
			if typeopr.IsEmpty(val) {
				if err := emptyOperations(&val, t.Field(i).Tag.Get(namelib.TAGS.DB_MAPPER_EMPTY), tag); err != nil {
					return err
				}
			}
			value.SetString(string(val))
		}
	}
	return nil
}

// ParamsValueFromStruct creates a map from a structure that describes the table.
// To work correctly, you need a completed structure, and the required fields must have the `name:"<column name>"` tag.
func ParamsValueFromStruct(filledStructurePtr typeopr.IPtr, nilIfEmpty []string) (map[string]any, error) {
	structure := filledStructurePtr.Ptr()
	if !typeopr.PtrIsStruct(structure) {
		return nil, typeopr.ErrParameterNotStruct{Param: "structure"}
	}
	outputParamsMap := make(map[string]any)

	typeof := reflect.TypeOf(structure).Elem()
	valueof := reflect.ValueOf(structure).Elem()
	for i := 0; i < typeof.NumField(); i++ {
		fieldValue := valueof.Field(i)
		dbColName := typeof.Field(i).Tag.Get(namelib.TAGS.DB_MAPPER_NAME)
		if dbColName == "" {
			continue
		}
		if fslice.SliceContains(nilIfEmpty, dbColName) && fieldValue.IsZero() {
			outputParamsMap[dbColName] = nil
		} else {
			outputParamsMap[dbColName] = fieldValue.Interface()
		}
	}
	return outputParamsMap, nil
}

func validateDbTag(tag string, dbRes *map[string]interface{}) (bool, error) {
	if tag == "" {
		return false, nil
	}
	if _, ok := (*dbRes)[tag]; !ok {
		return false, ErrDbResFieldNotFound{Field: tag}
	}
	return true, nil
}

func dbValueConversionToByte(value interface{}) ([]byte, error) {
	if value == nil {
		return nil, nil
	}
	var val []byte
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		v, ok := value.([]uint8)
		if !ok {
			return nil, typeopr.ErrConvertType{Type1: reflect.TypeOf(value).String(), Type2: "[]uint8"}
		}
		val = v
	} else {
		v, err := typeopr.AnyToBytes(value)
		if err != nil {
			return nil, err
		}
		val = v
	}
	return val, nil
}

func emptyOperations(val *[]byte, emptyTagValue string, dbFieldName string) error {
	if emptyTagValue == "" {
		return nil
	}
	switch emptyTagValue {
	case "-error":
		return typeopr.ErrValueIsEmpty{Value: dbFieldName}
	default:
		_val, err := typeopr.AnyToBytes(emptyTagValue)
		if err != nil {
			return err
		}
		*val = _val
	}
	return nil
}

type ErrDbResFieldNotFound struct {
	Field string
}

func (e ErrDbResFieldNotFound) Error() string {
	return fmt.Sprintf("The database output result does not contain the %s field.", e.Field)
}
