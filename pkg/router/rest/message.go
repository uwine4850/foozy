package rest

import "fmt"

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
