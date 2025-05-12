package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

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

var typeTime = reflect.TypeOf(time.Time{})
var typeBoolean = reflect.TypeOf(true)
var typeMap = reflect.TypeOf(map[string]interface{}{})
var typeByte = reflect.TypeOf([]byte{})
var typeDecimal = reflect.TypeOf(decimal.Decimal{})

func convertDBType(value *reflect.Value, structTag *reflect.StructTag, dataPtr *any) error {
	data := *dataPtr
	switch value.Type() {
	case typeTime:
		val, err := dbValueConversionToByte(data)
		if err != nil {
			return err
		}
		pTime, err := convertTimeF(&val, structTag.Get(namelib.TAGS.DB_MAPPER_DATE_F))
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(pTime))
	case typeBoolean:
		cVal, err := convertBoolean(dataPtr)
		if err != nil {
			return err
		}
		value.SetBool(cVal)
	case typeMap:
		v := map[string]interface{}{}
		if err := json.Unmarshal(data.([]byte), &v); err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	case typeByte:
		value.Set(reflect.ValueOf(data))
	case typeDecimal:
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

type ErrTypeConversion struct {
	Type string
}

func (e ErrTypeConversion) Error() string {
	return fmt.Sprintf("%s type conversion error", e.Type)
}
