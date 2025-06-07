package rest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

// ImplementDTOMessage structure to be embedded in a message.
// Once inlined, the framework will implement the irest.IMessage interface.
type ImplementDTOMessage struct{}

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

// IsSafeMessage checks whether the message is safe.
// A message is safe if it is in allowed messages.
func IsSafeMessage(message typeopr.IPtr, allowedMessages []AllowMessage) error {
	_type := reflect.TypeOf(message.Ptr()).Elem()
	if !typeopr.IsImplementInterface(message, (*irest.IMessage)(nil)) {
		return fmt.Errorf("%s message does not implement irest.IMessage interface", _type)
	}
	// If the message type is passed through the irest.IMessage interface.
	if _type == reflect.TypeOf((*irest.IMessage)(nil)).Elem() {
		_type = reflect.TypeOf(reflect.ValueOf(message.Ptr()).Elem().Interface())
	}
	pkgAndName := strings.Split(_type.String(), ".")
	msg := AllowMessage{Package: pkgAndName[0], Name: pkgAndName[1]}
	if !fslice.SliceContains(allowedMessages, msg) {
		return fmt.Errorf("%s message is unsafe", msg.FullName())
	}
	return nil
}
