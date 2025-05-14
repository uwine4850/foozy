package rest

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"sync"
	"text/template"

	"strings"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/form"
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

// DTO (Data Transfer Object) generates a typescript interface using a message.
// For proper operation, you must make sure that the allowed messages match the
// transferred messages.
//
// IMPORTANT: Allowed messages [AllowedMessages] must be exactly the same as the
// messages to generate [Messages].  This is only needed during genarration.
// If no generation is taking place, you can use the allowed [AllowedMessages] alone.
//
// Any dependencies must also be included in the allowed messages, and they must be
// in the same file as the parent object. Importing from other files is not allowed.
//
// Structure fields must have the `dto:“<field_name>”` tag for successful generation.
// Where <field_name> is the name of the field in the typescript interface.
// If the tag is not added, the field will simply be skipped during generation.
type DTO struct {
	allowedMessages []AllowMessage
	messages        map[string][]irest.IMessage
	isGenerated     bool
}

func NewDTO() *DTO {
	once.Do(func() {
		instance = &DTO{}
	})
	return instance
}

// AllowedMessages list of allowed messages for generation.
// These messages will be checked before using the message,
// if it is not in this list, an error will be raised.
func (d *DTO) AllowedMessages(messages []AllowMessage) {
	d.allowedMessages = messages
}

func (d *DTO) GetAllowedMessages() []AllowMessage {
	return d.allowedMessages
}

// Messages a list of messages that will be used for generation.
// So will check the generation allowances for each message.
func (d *DTO) Messages(messages map[string][]irest.IMessage) {
	d.messages = messages
}

// Generate start of generation.
// Uses a template to generate typescript interfaces.
func (d *DTO) Generate() error {
	if d.isGenerated {
		return ErrMultipleGenerateCall{}
	}
	allGeneratedAllowMessages := []AllowMessage{}
	acceptMessages := map[string][]genMessage{}
	for path, messages := range d.messages {
		generatedMessages, generetedAllowMessages, err := d.getGenMessages(messages)
		if err != nil {
			return err
		}
		allGeneratedAllowMessages = append(allGeneratedAllowMessages, generetedAllowMessages...)
		acceptMessages[path] = generatedMessages
	}
	if err := d.validateGeneratedMessage(allGeneratedAllowMessages); err != nil {
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

// getGenMessages generates typescript interfaces and stores them
// in the special structures [genMessage]. Each such structure contains data of one interface.
// Also returns [AllowMessage]. This structure contains data about one generated DTO message.
func (d *DTO) getGenMessages(messages []irest.IMessage) ([]genMessage, []AllowMessage, error) {
	generatedMessages := []genMessage{}
	generatedAllowMessages := []AllowMessage{}
	for i := 0; i < len(messages); i++ {
		_type := reflect.TypeOf((messages)[i])
		typeInfo := strings.Split(_type.String(), ".")

		allowedMessage := AllowMessage{Package: typeInfo[0], Name: typeInfo[1]}
		if !slices.Contains(d.allowedMessages, allowedMessage) {
			return nil, nil, ErrMessageNotAllowed{MessageType: allowedMessage.FullName()}
		}
		var genMsg genMessage
		genMsg.Name = _type.Name()
		for i := 0; i < _type.NumField(); i++ {
			// Skip implemetation object.
			if _type.Field(i).Type == reflect.TypeOf(ImplementDTOMessage{}) {
				continue
			}
			cnvType, err := d.convertType(_type.Field(i).Type, messages, allowedMessage)
			if err != nil {
				return nil, nil, err
			}
			if tagFieldName := _type.Field(i).Tag.Get(namelib.TAGS.REST_MAPPER_NAME); tagFieldName != "" {
				// Formation of the [genMessageField] structure.
				messageField := genMessageField{Type: cnvType}
				messageField.Name = tagFieldName
				genMsg.Fields = append(genMsg.Fields, messageField)
			} else {
				// Skip if tag no exists.
				continue
			}
		}
		if len(genMsg.Fields) == 0 {
			return nil, nil, ErrNumberOfFields{MessageType: allowedMessage.FullName()}
		}
		generatedMessages = append(generatedMessages, genMsg)
		generatedAllowMessages = append(generatedAllowMessages, allowedMessage)
	}
	return generatedMessages, generatedAllowMessages, nil
}

func (d *DTO) convertType(goType reflect.Type, messages []irest.IMessage, mainMessage AllowMessage) (string, error) {
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
		if goType == reflect.TypeOf(form.FormFile{}) {
			return "File | null", nil
		}
		typeInfo := strings.Split(goType.String(), ".")
		if !slices.Contains(d.allowedMessages, AllowMessage{Package: typeInfo[0], Name: typeInfo[1]}) {
			return "", ErrMessageNotAllowed{MessageType: goType.String()}
		}
		for i := 0; i < len(messages); i++ {
			if reflect.TypeOf((messages)[i]) == goType {
				return fmt.Sprintf("%s | undefined", goType.Name()), nil
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
		return fmt.Sprintf("Record<%s, %s>", cnvKeyType, cnvValueType), nil
	default:
		return "", fmt.Errorf("unsupported data type: %s", goType.String())
	}
}

// validateGeneratedMessage validation of generated messages.
// Checks if the generated messages match the allowed messages.
func (d *DTO) validateGeneratedMessage(generatedMessage []AllowMessage) error {
	for i := 0; i < len(d.allowedMessages); i++ {
		if !slices.Contains(generatedMessage, d.allowedMessages[i]) {
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

type ErrMultipleGenerateCall struct{}

func (e ErrMultipleGenerateCall) Error() string {
	return "generate can only be called once"
}
