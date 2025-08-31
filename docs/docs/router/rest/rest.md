## dto
The DTO system is an implementation of REST API with some special features. The purpose of DTO is to control and type the messages being sent.<br>
To make interaction with the client more convenient, you can use TypeScript interface generation.<br>
It works as follows:

1. A message (structure) is created that controls the type of data being received and sent.
2. All messages are registered.
3. __Optional__. TypeScript interfaces are generated that correspond to the server message (structure).
4. Messages are sent or received using special functions and methods.

Example of use in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/router_test/dto_test).

### DTO
DTO (Data Transfer Object) generates a typescript interface using a message.
For proper operation, you must make sure that the allowed messages match the transferred messages.

__IMPORTANT__: Allowed messages [AllowedMessages] must be exactly the same as the messages to generate [Messages]. This is only needed 
during genaration. If no generation is taking place, you can use the allowed [AllowedMessages] alone.

Any dependencies must also be included in the allowed messages, and they must be 
in the same file as the parent object. Importing from other files is not allowed.

Structure fields must have the `dto:“<field_name>”` tag for successful generation. Where `<field_name>` is the name of the field in 
the typescript interface. If the tag is not added, the field will simply be skipped during generation.
```golang
type DTO struct {
	allowedMessages []AllowMessage
	messages        map[string][]irest.Message
	isGenerated     bool
}
```

#### DTO.AllowedMessages
List of allowed messages for generation.<br>
These messages will be checked before using the message, if it is not in this list, an error will be raised.
```golang
func (d *DTO) AllowedMessages(messages []AllowMessage) {
	d.allowedMessages = messages
}
```

#### DTO.GetAllowedMessages
Return `[]AllowMessage`.
```golang
func (d *DTO) GetAllowedMessages() []AllowMessage {
	return d.allowedMessages
}
```

#### DTO.Messages
A list of messages that will be used for generation.<br>
So will check the generation allowances for each message.
```golang
func (d *DTO) Messages(messages map[string][]irest.Message) {
	d.messages = messages
}
```

#### DTO.Generate
Start of generation.<br>
Uses a template to generate typescript interfaces.
```golang
func (d *DTO) Generate() error {
	if d.isGenerated {
		return ErrMultipleGenerateCall{}
	}
	if err := d.validateMessageIntegrity(); err != nil {
		return err
	}
	if err := d.writeClientMessages(); err != nil {
		return err
	}
	d.isGenerated = true
	return nil
}
```

### ImplementDTOMessage
Structure to be embedded in a message.<br>
Once inlined, the framework will implement the irest.IMessage interface.
```golang
type ImplementDTOMessage struct{}
```

### AllowMessage
Used to transmit packet data and message name in string type.
```golang
type AllowMessage struct {
	Package string
	Name    string
}
```

#### AllowMessage.FullName
Outputs the full name of the message.<br>
For example: main.Message.
```golang
func (a *AllowMessage) FullName() string {
	return fmt.Sprintf("%s.%s", a.Package, a.Name)
}
```

#### IsSafeMessage
Checks whether the message is safe.<br>
A message is safe if it is in allowed messages.
```golang
func IsSafeMessage(message typeopr.IPtr, allowedMessages []AllowMessage) error {
	_type := reflect.TypeOf(message.Ptr()).Elem()
	if !typeopr.IsImplementInterface(message, (*irest.Message)(nil)) {
		return fmt.Errorf("%s message does not implement irest.IMessage interface", _type)
	}
	// If the message type is passed through the irest.IMessage interface.
	if _type == reflect.TypeOf((*irest.Message)(nil)).Elem() {
		_type = reflect.TypeOf(reflect.ValueOf(message.Ptr()).Elem().Interface())
	}
	pkgAndName := strings.Split(_type.String(), ".")
	msg := AllowMessage{Package: pkgAndName[0], Name: pkgAndName[1]}
	if !fslice.SliceContains(allowedMessages, msg) {
		return fmt.Errorf("%s message is unsafe", msg.FullName())
	}
	return nil
}
```