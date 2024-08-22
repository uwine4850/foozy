package restmapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

// FillMessageFromMap fills in a message from the card.
// To work you need to use the 'json' tag.
func FillMessageFromMap(jsonMap *map[string]interface{}, outputPtr itypeopr.IPtr) error {
	output := outputPtr.Ptr()
	if !typeopr.IsImplementInterface(typeopr.Ptr{}.New(output), (*irest.IMessage)(nil)) {
		return errors.New("output param must implement the irest.IMessage interface")
	}
	var v reflect.Value
	var t reflect.Type
	if reflect.DeepEqual(reflect.TypeOf(output), reflect.TypeOf(&reflect.Value{})) {
		v = *output.(*reflect.Value)
		t = v.Type()
	} else {
		t = reflect.TypeOf(output).Elem()
		v = reflect.ValueOf(output).Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		inputValue, ok := (*jsonMap)[name]
		if !ok {
			continue
		}
		fieldType := t.Field(i).Type
		fieldValue := v.Field(i)
		switch fieldType.Kind() {
		case reflect.Struct:
			vv := inputValue.(map[string]interface{})
			if err := FillMessageFromMap(&vv, typeopr.Ptr{}.New(&fieldValue)); err != nil {
				return err
			}
		default:
			if err := fillField(&fieldValue, inputValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func fillField(fieldValue *reflect.Value, fieldData interface{}) error {
	switch fieldValue.Type().Kind() {
	case reflect.Map:
		fieldDataValue := reflect.ValueOf(fieldData)
		newMap := reflect.MakeMap(fieldValue.Type())
		fieldDataValueMap := fieldDataValue.MapRange()
		for fieldDataValueMap.Next() {
			val := reflect.New(fieldValue.Type().Elem()).Elem()
			if err := fillField(&val, fieldDataValueMap.Value().Interface()); err != nil {
				return err
			}
			newMap.SetMapIndex(fieldDataValueMap.Key(), val)
		}
		fieldValue.Set(newMap)
	case reflect.Struct:
		fd := fieldData.(map[string]interface{})
		if err := FillMessageFromMap(&fd, typeopr.Ptr{}.New(fieldValue)); err != nil {
			return err
		}
	case reflect.Slice:
		dataValue := reflect.ValueOf(fieldData)
		newSlice := reflect.MakeSlice(fieldValue.Type(), 0, dataValue.Len())
		for i := 0; i < dataValue.Len(); i++ {
			val := reflect.New(fieldValue.Type().Elem()).Elem()
			if err := fillField(&val, dataValue.Index(i).Interface()); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, val)
		}
		fieldValue.Set(newSlice)
	case reflect.Int:
		// It needs to be converted to float64 because json.Unmarshal represents any numbers only in this format.
		fl := reflect.ValueOf(fieldData).Interface().(float64)
		fieldValue.Set(reflect.ValueOf(int(fl)))
	case reflect.String:
		fieldValue.Set(reflect.ValueOf(template.HTMLEscapeString(fieldData.(string))))
	case reflect.Float32, reflect.Float64:
		fieldValue.Set(reflect.ValueOf(fieldData))
	case reflect.Bool:
		fieldValue.Set(reflect.ValueOf(fieldData))
	default:
		return fmt.Errorf("converting JSON to %s type is not supported", fieldValue.Type().Kind().String())
	}
	return nil
}

// JsonStringToMap converts a JSON string into a map.
func JsonStringToMap(data string, m *map[string]interface{}) error {
	if err := json.Unmarshal([]byte(data), m); err != nil {
		return err
	}
	return nil
}
