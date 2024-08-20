package rest

import (
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

// InmplementDTOMessage structure to be embedded in a message.
// Once inlined, the framework will implement the irest.IMessage interface.
type InmplementDTOMessage struct {
}

func (m InmplementDTOMessage) IsImplementDTOMessage() {}

// AllowMessage used to transmit packet data and message name in string type.
type AllowMessage struct {
	Package string
	Name    string
}

func (a *AllowMessage) FullName() string {
	return fmt.Sprintf("%s.%s", a.Package, a.Name)
}

func DeepCheckSafeMessage(dto *DTO, messagePtr itypeopr.IPtr) error {
	if err := dto.IsSafeMessage(messagePtr); err != nil {
		return err
	}
	message := messagePtr.Ptr()
	_type := reflect.TypeOf(message).Elem()
	value := reflect.ValueOf(message).Elem()
	for i := 0; i < _type.NumField(); i++ {
		field := _type.Field(i)
		v := value.Field(i)
		if field.Type.Kind() == reflect.Struct && !reflect.DeepEqual(field.Type, reflect.TypeOf(InmplementDTOMessage{})) {
			if err := DeepCheckSafeMessage(dto, typeopr.Ptr{}.New(v.Addr().Interface())); err != nil {
				return err
			}
		}
	}
	return nil
}