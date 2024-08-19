package rest

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"text/template"

	"strings"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

var (
	instance *DTO
	once     sync.Once
)

// A template for generating a typescript interface using a structure.
const tmpl = `
{{- range .}}
export interface {{.Name}} {
    {{- range .Fields}}
    {{ .Name }}: {{ .Type }};
    {{- end}}
}
{{- end}}
`

// DTO(Data Transfer Object) generates a typescript interface using message.
// For proper operation, you need to make sure that the allowed messages match the transmitted ones.
// It is important to understand that allowed messages must ALL be used, that is, if a message is allowed, it must always be used when generated.
// Any dependencies also need to be included in resolved messages and they must be in the same file as the parent object.
type DTO struct {
	allowedMessages []AllowMessage
	messages        map[string]*[]irest.IMessage
	isGenerated     bool
}

func NewDTO() *DTO {
	once.Do(func() {
		instance = &DTO{}
	})
	return instance
}

// AllowedMessages list of allowed messages for generation.
func (d *DTO) AllowedMessages(messages []AllowMessage) {
	if d.allowedMessages != nil {
		panic("AllowedMessages can only be call once")
	}
	d.allowedMessages = messages
}

func (d *DTO) GetAllowedMessages() []AllowMessage {
	return d.allowedMessages
}

// Messages a list of messages that will be used for generation.
func (d *DTO) Messages(messages map[string]*[]irest.IMessage) {
	if d.messages != nil {
		panic("Messages can only be call once")
	}
	d.messages = messages
}

// Generate start of generation.
func (d *DTO) Generate() error {
	if d.isGenerated {
		panic("Generate can only be call once")
	}
	allGeneretedAllowMessages := []AllowMessage{}
	acceptMessages := map[string][]genMessage{}
	for path, messages := range d.messages {
		generetedMessaages, generetedAllowMessages, err := d.getGenMessaages(messages)
		if err != nil {
			return err
		}
		allGeneretedAllowMessages = append(allGeneretedAllowMessages, generetedAllowMessages...)
		acceptMessages[path] = generetedMessaages
	}
	if err := d.validateGeneratedMessage(allGeneretedAllowMessages); err != nil {
		return err
	}
	for path, genMessages := range acceptMessages {
		t := template.Must(template.New("code").Parse(tmpl))
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := t.Execute(file, genMessages); err != nil {
			return err
		}
	}
	d.isGenerated = true
	return nil
}

func (d *DTO) IsSafeMessage(message itypeopr.IPtr) error {
	_type := reflect.TypeOf(message.Ptr()).Elem()
	if !typeopr.IsImplementInterface(message, (*irest.IMessage)(nil)) {
		return fmt.Errorf("%s message does not implement irest.IMessage interface", _type)
	}
	typeInfo := strings.Split(_type.String(), ".")
	msg := AllowMessage{Package: typeInfo[0], Name: typeInfo[1]}
	if !fslice.SliceContains(d.allowedMessages, msg) {
		return fmt.Errorf("%s message is unsafe", msg.FullName())
	}
	return nil
}

func (d *DTO) getGenMessaages(messages *[]irest.IMessage) ([]genMessage, []AllowMessage, error) {
	generatedMessages := []genMessage{}
	generatedAllowMessages := []AllowMessage{}
	for i := 0; i < len(*messages); i++ {
		_type := reflect.TypeOf((*messages)[i])
		typeInfo := strings.Split(_type.String(), ".")

		allowMessage := AllowMessage{Package: typeInfo[0], Name: typeInfo[1]}
		if !fslice.SliceContains(d.allowedMessages, allowMessage) {
			return nil, nil, ErrMessageNotAllowed{MessageType: allowMessage.FullName()}
		}
		var genMsg genMessage
		genMsg.Name = _type.Name()
		for i := 0; i < _type.NumField(); i++ {
			if _type.Field(i).Type == reflect.TypeOf(InmplementDTOMessage{}) {
				continue
			}
			cnvType, err := d.convertType(_type.Field(i).Type, messages, allowMessage)
			if err != nil {
				return nil, nil, err
			}
			messageField := genMessageField{Name: _type.Field(i).Name, Type: cnvType}
			genMsg.Fields = append(genMsg.Fields, messageField)
		}
		if len(genMsg.Fields) == 0 {
			return nil, nil, ErrNumberOfFields{MessageType: allowMessage.FullName()}
		}
		generatedMessages = append(generatedMessages, genMsg)
		generatedAllowMessages = append(generatedAllowMessages, allowMessage)
	}
	return generatedMessages, generatedAllowMessages, nil
}

func (d *DTO) convertType(goType reflect.Type, messages *[]irest.IMessage, mainMessage AllowMessage) (string, error) {
	switch goType.Kind() {
	case reflect.Int, reflect.Float64, reflect.Float32:
		return "number", nil
	case reflect.String:
		return "string", nil
	case reflect.Bool:
		return "boolean", nil
	case reflect.Slice:
		cnvType, err := d.convertType(goType.Elem(), messages, mainMessage)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s[]", cnvType), nil
	case reflect.Struct:
		typeInfo := strings.Split(goType.String(), ".")
		if !fslice.SliceContains(d.allowedMessages, AllowMessage{Package: typeInfo[0], Name: typeInfo[1]}) {
			return "", ErrMessageNotAllowed{MessageType: goType.String()}
		}
		for i := 0; i < len(*messages); i++ {
			if reflect.TypeOf((*messages)[i]) == goType {
				return goType.Name(), nil
			}
		}
		return "", ErrNoDependency{DependencyType: goType.String(), MessageType: mainMessage.FullName()}
	case reflect.Map:
		cnvKeyType, err := d.convertType(goType.Key(), messages, mainMessage)
		if err != nil {
			return "", err
		}
		cnvValueType, err := d.convertType(goType.Elem(), messages, mainMessage)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Map<%s, %s>", cnvKeyType, cnvValueType), nil
	default:
		panic(fmt.Sprintf("%s data type is not supported", goType))
	}
}

func (d *DTO) validateGeneratedMessage(generatedMessage []AllowMessage) error {
	for i := 0; i < len(d.allowedMessages); i++ {
		if !fslice.SliceContains(generatedMessage, d.allowedMessages[i]) {
			return ErrMessageNotImplemented{MessageType: d.allowedMessages[i].FullName()}
		}
	}
	return nil
}

type genMessage struct {
	Name   string
	Fields []genMessageField
}

type genMessageField struct {
	Name string
	Type string
}

type ErrMessageNotAllowed struct {
	MessageType string
}

func (e ErrMessageNotAllowed) Error() string {
	return fmt.Sprintf("message %s is not allowed", e.MessageType)
}

type ErrNumberOfFields struct {
	MessageType string
}

func (e ErrNumberOfFields) Error() string {
	return fmt.Sprintf("the number of message fields %s must be greater than 0", e.MessageType)
}

type ErrNoDependency struct {
	DependencyType string
	MessageType    string
}

func (e ErrNoDependency) Error() string {
	return fmt.Sprintf("dependency %s not found for message %s", e.DependencyType, e.MessageType)
}

type ErrMessageNotImplemented struct {
	MessageType string
}

func (e ErrMessageNotImplemented) Error() string {
	return fmt.Sprintf("%s message not implemented", e.MessageType)
}
