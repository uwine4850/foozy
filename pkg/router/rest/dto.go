package rest

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"sync"
	"text/template"

	"strings"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
export function is{{ .Name }}(obj: any): obj is {{ .Name }} {
    return typeof obj === 'object' && obj !== null && '{{ .TypeId }}' in obj;
}
{{- end}}
`

const serverMessageTemplate = `package {{ .PkgName }}
import "github.com/uwine4850/foozy/pkg/router/rest"
{{ range $import := .Imports }}
{{ $import }}
{{- end }}

{{- range .Messages}}

type {{.Name}} struct {
    rest.ImplementDTOMessage
    {{ .MessageID }}
    {{- range .Fields}}
    {{ .Name }} {{ .Type }}
    {{- end}}
}

func {{ .FuncName }}({{ range $name, $type := .FuncArgs }}
    {{ $name | ToLower }} {{ $type }},
{{- end }}
) *{{.Name}} {
    return &{{.Name}}{ {{ range $name, $type := .FuncArgs }}
	    {{$name}}: {{$name | ToLower}},
{{- end }}
    }
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
	if err := d.validateMessageIntegrity(); err != nil {
		return err
	}
	if err := d.writeClientMessages(); err != nil {
		return err
	}
	if err := d.writeServerMessages(); err != nil {
		return err
	}
	d.isGenerated = true
	return nil
}

// validateMessageIntegrity checks if the allowed messages match the passed messages.
func (d *DTO) validateMessageIntegrity() error {
	for _, messages := range d.messages {
		for i := 0; i < len(messages); i++ {
			pkgAndName := strings.Split(reflect.TypeOf(messages[i]).String(), ".")
			pkgName := pkgAndName[0]
			messageName := pkgAndName[1]
			allowMessage := AllowMessage{Package: pkgName, Name: messageName}
			if !slices.Contains(d.allowedMessages, allowMessage) {
				return ErrMessageNotAllowed{MessageType: allowMessage.FullName()}
			}
		}
	}
	return nil
}

// writeClientMessages writes messages to client files that are passed through the [Messages] method.
func (d *DTO) writeClientMessages() error {
	allGeneratedAllowMessages := []AllowMessage{}
	acceptMessages := map[string][]clientMessage{}
	for path, messages := range d.messages {
		generatedMessages, generetedAllowMessages, err := d.generateMessages(messages)
		if err != nil {
			return err
		}
		allGeneratedAllowMessages = append(allGeneratedAllowMessages, generetedAllowMessages...)
		acceptMessages[path] = generatedMessages
	}
	if err := d.validateImplementMessage(allGeneratedAllowMessages); err != nil {
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
	return nil
}

// writeServerMessages writes generated server message data to a file.
// To work it is necessary to specify the path to the file in the DTO configuration.
func (d *DTO) writeServerMessages() error {
	serverDTOPath := config.LoadedConfig().Default.DTOConfig.GeneratedFilePath
	if serverDTOPath == "" {
		return errors.New("the Default.DTOConfig.GeneratedFilePath configuration cannot be empty")
	}
	serverDTOFile, err := os.Create(serverDTOPath)
	if err != nil {
		return err
	}
	defer serverDTOFile.Close()

	serverFile, err := d.generateServerFile()
	if err != nil {
		return err
	}
	funcMap := template.FuncMap{
		"ToLower": fstring.ToLower,
	}
	st := template.Must(template.New("serverM").Funcs(funcMap).Parse(serverMessageTemplate))
	if err := st.Execute(serverDTOFile, serverFile); err != nil {
		return err
	}
	return nil
}

// generateServerFile generates a single file for all server messages.
// It is important to specify the package name and file path in the settings.
func (d *DTO) generateServerFile() (*generatedServerFile, error) {
	newServerMessages := []serverMessage{}
	imports := []string{}
	for _, messages := range d.messages {
		for msgID := 0; msgID < len(messages); msgID++ {
			messageType := reflect.TypeOf(messages[msgID])
			newServerMessage := serverMessage{}
			pkgName := strings.Split(messageType.String(), ".")[0]
			newServerMessage.Name = pkgName + "_" + messageType.Name()
			funcArgs := map[string]string{}
			for fieldID := 0; fieldID < messageType.NumField(); fieldID++ {
				field := messageType.Field(fieldID)
				dtoTag := field.Tag.Get(namelib.TAGS.REST_MAPPER_NAME)
				if field.Type != reflect.TypeOf(ImplementDTOMessage{}) && dtoTag != "" {
					templateFieldType, templateArgFieldType := getTemplateFieldAndArgFieldType(field)
					if field.Type == reflect.TypeOf(form.FormFile{}) {
						imports = append(imports, `import "github.com/uwine4850/foozy/pkg/router/form"`)
						templateFieldType = fmt.Sprintf("form.%s `%s:\"%s\"`", field.Type.Name(), namelib.TAGS.REST_MAPPER_NAME, dtoTag)
						templateArgFieldType = "form." + field.Type.Name()
					}
					newGeneratedMessageField := generatedMessageField{
						Name: field.Name,
						Type: templateFieldType,
					}
					newServerMessage.Fields = append(newServerMessage.Fields, newGeneratedMessageField)
					funcArgs[field.Name] = templateArgFieldType
				}
			}
			pkgAndStructName := strings.Split(messageType.String(), ".")
			dtoMessageIdName := fmt.Sprintf("Type%s", strings.Replace(messageType.String(), ".", "", -1))
			dtoMessageIdTypeName := fmt.Sprintf("any `%s:\"%s\"`",
				namelib.TAGS.REST_MAPPER_NAME, dtoMessageIdName)
			newServerMessage.MessageID = fmt.Sprintf("%s %s", dtoMessageIdName, dtoMessageIdTypeName)
			funcName := cases.Title(language.Und).String(pkgAndStructName[0]) + pkgAndStructName[1]
			newServerMessage.FuncName = fmt.Sprintf("New%s", funcName)
			newServerMessage.FuncArgs = funcArgs
			newServerMessages = append(newServerMessages, newServerMessage)
		}
	}
	pkgName := config.LoadedConfig().Default.DTOConfig.PkgName
	if pkgName == "" {
		return nil, errors.New("package name cannot be empty")
	}
	return &generatedServerFile{
		PkgName:  pkgName,
		Messages: newServerMessages,
		Imports:  imports,
	}, nil
}

// generateMessages generates typescript interfaces and stores them
// in the special structures [genMessage]. Each such structure contains data of one interface.
// Also returns [AllowMessage]. This structure contains data about one generated DTO message.
func (d *DTO) generateMessages(messages []irest.IMessage) ([]clientMessage, []AllowMessage, error) {
	clientMessages := []clientMessage{}
	generatedAllowMessages := []AllowMessage{}
	for i := 0; i < len(messages); i++ {
		_type := reflect.TypeOf((messages)[i])
		typeInfo := strings.Split(_type.String(), ".")

		allowedMessage := AllowMessage{Package: typeInfo[0], Name: typeInfo[1]}
		var clientMessage clientMessage
		clientMessage.Name = _type.Name()
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
				messageField := generatedMessageField{Type: cnvType, Name: tagFieldName}
				clientMessage.Fields = append(clientMessage.Fields, messageField)
			} else {
				// Skip if tag no exists.
				continue
			}
		}
		if len(clientMessage.Fields) == 0 {
			return nil, nil, ErrNumberOfFields{MessageType: allowedMessage.FullName()}
		}

		typeId := fmt.Sprintf("Type%s", strings.Replace(allowedMessage.FullName(), ".", "", -1))
		typeMessageFiels := generatedMessageField{
			Name: typeId + "?",
			Type: "unknown",
		}
		clientMessage.Fields = append([]generatedMessageField{typeMessageFiels}, clientMessage.Fields...)
		clientMessage.TypeId = typeId
		clientMessages = append(clientMessages, clientMessage)
		generatedAllowMessages = append(generatedAllowMessages, allowedMessage)
	}
	return clientMessages, generatedAllowMessages, nil
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
		// This check is needed to make sure that the structure is safe and is a valid message and not any other structure.
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

// validateImplementMessage validation of generated messages.
// Checks if the generated messages match the allowed messages.
func (d *DTO) validateImplementMessage(generatedMessage []AllowMessage) error {
	for i := 0; i < len(d.allowedMessages); i++ {
		if !slices.Contains(generatedMessage, d.allowedMessages[i]) {
			return ErrMessageNotImplemented{MessageType: d.allowedMessages[i].FullName()}
		}
	}
	return nil
}

type clientMessage struct {
	Name   string
	TypeId string
	Fields []generatedMessageField
}

type generatedServerFile struct {
	PkgName  string
	Imports  []string
	Messages []serverMessage
}

type serverMessage struct {
	PkgName   string
	Name      string
	MessageID string
	Fields    []generatedMessageField
	FuncName  string
	// <name>:<type>
	FuncArgs map[string]string
}

type generatedMessageField struct {
	Name string
	Type string
}

func getTemplateFieldAndArgFieldType(field reflect.StructField) (templFieldType string, templArgFieldType string) {
	dtoTag := field.Tag.Get(namelib.TAGS.REST_MAPPER_NAME)
	var templateFieldType string
	var templateArgFieldType string
	if field.Type.Kind() == reflect.Struct {
		fieldStructPkg := strings.Split(field.Type.String(), ".")[0]
		fieldStructTypeName := fieldStructPkg + "_" + field.Type.Name()
		templateFieldType = fmt.Sprintf("%s `%s:\"%s\"`", fieldStructTypeName, namelib.TAGS.REST_MAPPER_NAME, dtoTag)
		templateArgFieldType = fieldStructTypeName
	} else {
		templateFieldType = fmt.Sprintf("%s `%s:\"%s\"`", field.Type.Name(), namelib.TAGS.REST_MAPPER_NAME, dtoTag)
		templateArgFieldType = field.Type.Name()
	}
	return templateFieldType, templateArgFieldType
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
