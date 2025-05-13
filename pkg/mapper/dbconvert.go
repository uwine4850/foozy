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

type DatabaseConverter struct{}

func (dc DatabaseConverter) dbValueConversionToByte(value interface{}) ([]byte, error) {
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

var (
	typeTime    = reflect.TypeOf(time.Time{})
	typeBool    = reflect.TypeOf(true)
	typeMap     = reflect.TypeOf(map[string]interface{}{})
	typeBytes   = reflect.TypeOf([]byte{})
	typeDecimal = reflect.TypeOf(decimal.Decimal{})
)

func (dc DatabaseConverter) convertDBType(value *reflect.Value, structTag *reflect.StructTag, dataPtr *any) error {
	data := *dataPtr
	switch value.Type() {
	case typeTime:
		val, err := dc.dbValueConversionToByte(data)
		if err != nil {
			return err
		}
		pTime, err := dc.convertTimeF(&val, structTag.Get(namelib.TAGS.DB_MAPPER_DATE_F))
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(pTime))
	case typeBool:
		cVal, err := dc.convertBoolean(dataPtr)
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
	case typeBytes:
		value.Set(reflect.ValueOf(data))
	case typeDecimal:
		cVal, err := dc.convertDecimal(dataPtr)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(cVal))
	default:
		if err := dc.convertOther(value, dataPtr); err != nil {
			return err
		}
	}
	return nil
}

func (dc DatabaseConverter) convertTimeF(val *[]byte, dateFormat string) (time.Time, error) {
	if dateFormat != "" {
		parsedTime, err := time.Parse(dateFormat, string(*val))
		if err != nil {
			return time.Time{}, err
		}
		return parsedTime, nil
	} else {
		formats := []string{time.DateTime, time.DateOnly, time.TimeOnly}
		for _, format := range formats {
			p, err := time.Parse(format, string(*val))
			if err == nil {
				return p, nil
			}
		}
		return time.Time{}, errors.New("invalid time format")
	}
}

func (dc DatabaseConverter) convertBoolean(dataPtr *any) (bool, error) {
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

func (dc DatabaseConverter) convertDecimal(dataPtr *any) (decimal.Decimal, error) {
	data := *dataPtr
	decimalString := string(data.([]byte))
	newDecimal, err := decimal.NewFromString(decimalString)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return newDecimal, nil
}

func (dc DatabaseConverter) convertOther(fieldValue *reflect.Value, dataPtr *any) error {
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
