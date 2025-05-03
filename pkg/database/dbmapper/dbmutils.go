package dbmapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

// FillStructFromDb Fills the structure with data from the database.
// For proper operation, you need to add a "db" tag with the column name for each field that represents a column from the database.
// For example: `db:"col_name"`.
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
		dbResField, ok := dbRes[tag]
		if ok {
			// Processing DB_MAPPER_EMPTY tag.
			if dbResField == nil {
				emptyVal := field.Tag.Get(namelib.TAGS.DB_MAPPER_EMPTY)
				if emptyVal != "" {
					if emptyVal == "-error" {
						return typeopr.ErrValueIsEmpty{Value: tag}
					}
					newByteData, err := dbValueConversionToByte(emptyVal)
					if err != nil {
						return err
					}
					newData := reflect.ValueOf(newByteData).Interface()
					if err := convertDBType(&value, &field, &newData); err != nil {
						return err
					}
					continue
				}
			}
			if err := convertDBType(&value, &field, &dbResField); err != nil {
				return err
			}
		}
	}
	return nil
}

// FillReflectValueFromDb fills structure of type reflect.Value with data from a database query.
// For proper operation, you need to add a "db" tag with the column name for each field that represents a column from the database.
// For example: `db:"col_name"`.
// You can also use the "empty" tag to indicate the action when the field is empty. If this tag is empty, nothing will happen. Other meanings:
// Any text - sets the field values ​​to this text.
// -error - display an error if the field is empty.
func FillReflectValueFromDb(dbRes map[string]interface{}, fill *reflect.Value) error {
	t := fill.Type()
	v := *fill
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := t.Field(i).Tag.Get(namelib.TAGS.DB_MAPPER_NAME)
		if exist, err := validateDbTag(tag, &dbRes); err != nil {
			return err
		} else {
			if !exist {
				continue
			}
		}
		dbResField, ok := dbRes[tag]
		if ok {
			// Processing DB_MAPPER_EMPTY tag.
			if dbResField == nil {
				emptyVal := field.Tag.Get(namelib.TAGS.DB_MAPPER_EMPTY)
				if emptyVal != "" {
					if emptyVal == "-error" {
						return typeopr.ErrValueIsEmpty{Value: tag}
					}
					newByteData, err := dbValueConversionToByte(emptyVal)
					if err != nil {
						return err
					}
					newData := reflect.ValueOf(newByteData).Interface()
					if err := convertDBType(&value, &field, &newData); err != nil {
						return err
					}
					continue
				}
			}
			if err := convertDBType(&value, &field, &dbResField); err != nil {
				return err
			}
		}
	}
	return nil
}

func convertDBType(value *reflect.Value, field *reflect.StructField, dataPtr *any) error {
	data := *dataPtr
	switch value.Type() {
	case reflect.TypeOf(time.Time{}):
		val, err := dbValueConversionToByte(data)
		if err != nil {
			return err
		}
		pTime, err := convertTimeF(&val, field.Tag.Get(namelib.TAGS.DB_MAPPER_DATE_F))
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(pTime))
	case reflect.TypeOf(true):
		cVal, err := convertBoolean(dataPtr)
		if err != nil {
			return err
		}
		value.SetBool(cVal)
	case reflect.TypeOf(map[string]interface{}{}):
		v := map[string]interface{}{}
		if err := json.Unmarshal(data.([]byte), &v); err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	case reflect.TypeOf([]byte{}):
		value.Set(reflect.ValueOf(data))
	case reflect.TypeOf(decimal.Decimal{}):
		cVal, err := convertDecimal(dataPtr)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(cVal))
	default:
		if err := convertOther(value, dataPtr); err != nil {
			return err
		}
	}
	return nil
}

func convertTimeF(val *[]byte, dateFormat string) (time.Time, error) {
	var pTime time.Time
	if dateFormat != "" {
		parsedTime, err := time.Parse(dateFormat, string(*val))
		if err != nil {
			return time.Time{}, err
		}
		pTime = parsedTime
	} else {
		var dateTimeError bool
		var dateError bool

		var parseError *time.ParseError

		// Datetime.
		parsedTime, err := time.Parse(time.DateTime, string(*val))
		if err != nil {
			if errors.As(err, &parseError) {
				dateTimeError = true
			} else {
				return time.Time{}, err
			}
		} else {
			pTime = parsedTime
		}
		if dateTimeError {
			// Date.
			parsedTime, err = time.Parse(time.DateOnly, string(*val))
			if err != nil {
				if errors.As(err, &parseError) {
					dateError = true
				} else {
					return time.Time{}, err
				}
			} else {
				pTime = parsedTime
			}
		}
		if dateError {
			// Time.
			parsedTime, err = time.Parse(time.TimeOnly, string(*val))
			if err != nil {
				return time.Time{}, err
			} else {
				pTime = parsedTime
			}
		}
	}
	return pTime, nil
}

func convertBoolean(dataPtr *any) (bool, error) {
	data := *dataPtr
	var zero int64
	if reflect.ValueOf(data).CanConvert(reflect.TypeOf(zero)) {
		int64Value := data.(int64)
		boolValue := int64Value != 0
		return boolValue, nil
	} else {
		return false, &ErrTypeConversion{Type: "boolean"}
	}
}

func convertDecimal(dataPtr *any) (decimal.Decimal, error) {
	data := *dataPtr
	decimalString := string(data.([]byte))
	newDecimal, err := decimal.NewFromString(decimalString)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return newDecimal, nil
}

func convertOther(fieldValue *reflect.Value, dataPtr *any) error {
	data := *dataPtr
	value := *fieldValue
	if value.Type().Kind() == reflect.Interface {
		value.Set(reflect.ValueOf(data))
	} else {
		if data == nil {
			return nil
		}
		if reflect.ValueOf(data).CanConvert(value.Type()) {
			convertValue := reflect.ValueOf(data).Convert(value.Type())
			value.Set(convertValue)
		} else {
			return ErrTypeConversion{Type: value.Type().Kind().String()}
		}
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

type ErrDbResFieldNotFound struct {
	Field string
}

func (e ErrDbResFieldNotFound) Error() string {
	return fmt.Sprintf("The database output result does not contain the %s field.", e.Field)
}

type ErrTypeConversion struct {
	Type string
}

func (e ErrTypeConversion) Error() string {
	return fmt.Sprintf("%s type conversion error", e.Type)
}
