package mapper

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"sync"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

var messageRawCache sync.Map
var implementDTOMessageType = reflect.TypeOf(rest.ImplementDTOMessage{})

// DeepCheckDTOSafeMessage checks whether transmitted messages and internal messages are safe.
// That is, there will be a check of internal structures, in depth to the limit.
// It is mandatory to have “dto” tags for each field.
//
// IMPORTANT: messagePtr accepts a pointer to a structure (preferably) or a pointer to
// a structure interface. In both cases, the object must implement the irest.IMessage interface.
func DeepCheckDTOSafeMessage(dto *rest.DTO, messagePtr typeopr.IPtr) error {
	if err := rest.IsSafeMessage(messagePtr, dto.GetAllowedMessages()); err != nil {
		return err
	}

	message := reflect.ValueOf(messagePtr.Ptr()).Elem()
	var RV reflect.Value
	if message.Type().Kind() == reflect.Interface {
		RV = message.Elem()
	} else {
		RV = message
	}
	rawObject := LoadSomeRawObjectFromCache(RV, &messageRawCache, namelib.TAGS.REST_MAPPER_NAME)
	for _, f := range *rawObject.Fields() {
		if f.Type.Kind() == reflect.Struct && f.Type != implementDTOMessageType {
			v := RV.FieldByName(f.Name)
			i := v.Interface().(irest.IMessage)
			if err := DeepCheckDTOSafeMessage(dto, typeopr.Ptr{}.New(&i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// JsonToMessage converts JSON data into the selected message.
// It is important that the message is safe.
func JsonToDTOMessage[T any](jsonData map[string]interface{}, dto *rest.DTO, output *T) error {
	if err := DeepCheckDTOSafeMessage(dto, typeopr.Ptr{}.New(output)); err != nil {
		return err
	}
	if err := FillDTOMessageFromMap(jsonData, output); err != nil {
		return err
	}
	return nil
}

// SendSafeJsonMessage sends only safe messages in JSON format.
func SendSafeJsonDTOMessage(w http.ResponseWriter, dto *rest.DTO, message typeopr.IPtr) error {
	if err := DeepCheckDTOSafeMessage(dto, message); err != nil {
		return err
	}
	if err := router.SendJson(message.Ptr(), w); err != nil {
		return err
	}
	return nil
}

// FillDTOMessageFromMap fills in a message from the card.
// To work you need to use the "dto" tag.
// If the DTO message is initially created correctly,
// there should be no problem with this function.
func FillDTOMessageFromMap[T any](jsonMap map[string]interface{}, out *T) error {
	if jsonMap == nil || out == nil {
		return errors.New("nil input to FillMessageFromMap")
	}
	RV := typeopr.GetReflectValue(out)
	if !typeopr.IsImplementInterface(typeopr.Ptr{}.New(out), (*irest.IMessage)(nil)) {
		return errors.New("output param must implement the irest.IMessage interface")
	}
	rawObject := LoadSomeRawObjectFromCache(RV, &messageRawCache, namelib.TAGS.REST_MAPPER_NAME)
	for name, f := range *rawObject.Fields() {
		inputValue, ok := (jsonMap)[name]
		if !ok {
			continue
		}
		fieldValue := RV.FieldByName(f.Name)
		switch f.Type.Kind() {
		case reflect.Struct:
			v, ok := inputValue.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected object for field '%s'", name)
			}
			if err := FillDTOMessageFromMap(v, &fieldValue); err != nil {
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
	fieldDataValue := reflect.ValueOf(fieldData)
	switch fieldValue.Type().Kind() {
	case reflect.Map:
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
		if err := FillDTOMessageFromMap(fd, fieldValue); err != nil {
			return err
		}
	case reflect.Slice:
		newSlice := reflect.MakeSlice(fieldValue.Type(), 0, fieldDataValue.Len())
		for i := 0; i < fieldDataValue.Len(); i++ {
			val := reflect.New(fieldValue.Type().Elem()).Elem()
			if err := fillField(&val, fieldDataValue.Index(i).Interface()); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, val)
		}
		fieldValue.Set(newSlice)
	case reflect.Int:
		// It needs to be converted to float64 because json.Unmarshal represents any numbers only in this format.
		fl := fieldDataValue.Interface().(float64)
		fieldValue.Set(reflect.ValueOf(int(fl)))
	case reflect.String:
		fieldValue.Set(reflect.ValueOf(template.HTMLEscapeString(fieldData.(string))))
	case reflect.Float32, reflect.Float64:
		fieldValue.Set(fieldDataValue)
	case reflect.Bool:
		fieldValue.Set(fieldDataValue)
	default:
		return fmt.Errorf("converting JSON to %s type is not supported", fieldValue.Type().Kind().String())
	}
	return nil
}
