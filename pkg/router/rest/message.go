package rest

import (
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

// ImplementDTOMessage structure to be embedded in a message.
// Once inlined, the framework will implement the irest.IMessage interface.
type ImplementDTOMessage struct {
}

func (m ImplementDTOMessage) IsImplementDTOMessage() {}

// AllowMessage used to transmit packet data and message name in string type.
type AllowMessage struct {
	Package string
	Name    string
}

// FullName outputs the full name of the message.
// For example: main.Message.
func (a *AllowMessage) FullName() string {
	return fmt.Sprintf("%s.%s", a.Package, a.Name)
}

// DeepCheckSafeMessage checks whether transmitted messages and internal messages are safe.
// That is, there will be a check of internal structures, in depth to the limit.
func DeepCheckSafeMessage(dto *DTO, messagePtr itypeopr.IPtr) error {
	if err := dto.IsSafeMessage(messagePtr); err != nil {
		return err
	}
	message := messagePtr.Ptr()
	_type := reflect.TypeOf(message).Elem()
	value := reflect.ValueOf(message).Elem()
	// If the message type is passed through the irest.IMessage interface.
	if _type == reflect.TypeOf((*irest.IMessage)(nil)).Elem() {
		_type = reflect.TypeOf(reflect.ValueOf(message).Elem())
		value = reflect.ValueOf(value.Interface())
	}
	for i := 0; i < _type.NumField(); i++ {
		field := _type.Field(i)
		v := value.Field(i)
		if field.Type.Kind() == reflect.Struct && !reflect.DeepEqual(field.Type, reflect.TypeOf(ImplementDTOMessage{})) {
			if err := DeepCheckSafeMessage(dto, typeopr.Ptr{}.New(v.Addr().Interface())); err != nil {
				return err
			}
		}
	}
	return nil
}
