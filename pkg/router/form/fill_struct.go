package form

import (
	"errors"
	"fmt"
	"github.com/uwine4850/foozy/pkg/ferrors"
	"github.com/uwine4850/foozy/pkg/utils"
	"mime/multipart"
	netUrl "net/url"
	"reflect"
)

type FormFile struct {
	Header *multipart.FileHeader
}

// FillStructFromForm A method that fills the structure with data from the form.
// The structure should always be passed as a pointer.
// For correct work it is necessary to specify "form" tag for each field of the structure. For example, `form:<form field name>`.
// Structure fields can be of two types only:
// []FormFile - form files.
// []string - all other data.
// The nilIfNotExist parameter sets the name of form fields that should be nil if they are not found (e.g. useful for checkboxes).
func FillStructFromForm(frm *Form, fill interface{}, nilIfNotExist []string) error {
	if reflect.TypeOf(fill).Kind() != reflect.Ptr {
		return ferrors.ErrParameterNotPointer{Param: "fill"}
	}
	if reflect.TypeOf(fill).Elem().Kind() != reflect.Struct {
		return ferrors.ErrParameterNotStruct{Param: "fill"}
	}
	formMap := FrmValueToMap(frm)
	t := reflect.TypeOf(fill).Elem()
	v := reflect.ValueOf(fill).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("form")
		// Skip if the tag is not a form
		if tag == "" {
			continue
		}
		err := checkFieldType(field)
		if err != nil {
			return err
		}
		formMapItem, ok := formMap[tag]
		if !ok {
			// Skips loop iteration if the field is not found, but it must be as nil.
			if utils.SliceContains(nilIfNotExist, tag) {
				continue
			} else {
				return ErrFormConvertFieldNotFound{tag}
			}
		}
		// Set files
		if reflect.DeepEqual(field.Type, reflect.TypeOf([]FormFile{})) && reflect.TypeOf(formMapItem) == reflect.TypeOf([]FormFile{}) {
			formType, _ := formMapItem.([]FormFile)
			if !ok {
				return ferrors.ErrConvertType{Type1: reflect.TypeOf(formMapItem).String(), Type2: "[]FormFile"}
			}
			value.Set(reflect.ValueOf(formType))
		}
		// Set string
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
			formType, ok := formMapItem.([]string)
			if !ok {
				return ferrors.ErrConvertType{Type1: reflect.TypeOf(formMapItem).String(), Type2: "string"}
			}
			value.Set(reflect.ValueOf(formType))
		}
	}
	return nil
}

// checkFieldType checks if the type of field to be filled is correct.
func checkFieldType(field reflect.StructField) error {
	if field.Type.Kind() != reflect.Slice {
		return errors.New(fmt.Sprintf("the %s field should be slice", field.Name))
	}
	if field.Type.Elem() != reflect.TypeOf("") && field.Type.Elem() != reflect.TypeOf(FormFile{}) {
		return errors.New(fmt.Sprintf("the %s field can only be of two types: []string or []FormFile.", field.Name))
	}
	return nil
}

type IFormGetEnctypeData interface {
	GetMultipartForm() *multipart.Form
	GetApplicationForm() netUrl.Values
}

// FrmValueToMap Converts the form to a map.
func FrmValueToMap(frm IFormGetEnctypeData) map[string]interface{} {
	formMap := make(map[string]interface{})
	multipartForm := frm.GetMultipartForm()
	if multipartForm != nil {
		for name, value := range multipartForm.Value {
			if value[0] == "" {
				formMap[name] = []string(nil)
			} else {
				formMap[name] = value
			}
		}
		for name, value := range multipartForm.File {
			var files []FormFile
			for i := 0; i < len(value); i++ {
				files = append(files, FormFile{Header: value[i]})
			}
			formMap[name] = files
		}
	}
	applicationForm := frm.GetApplicationForm()
	if applicationForm != nil {
		for name, value := range applicationForm {
			formMap[name] = value
		}
	}
	return formMap
}

type ErrArgumentNotPointer struct {
	Name string
}

func (e ErrArgumentNotPointer) Error() string {
	return fmt.Sprintf("argument %s is not a pointer", e.Name)
}

// FieldsNotEmpty checks the specified fields of the structure for emptiness.
// fieldsName - slice with exact names of structure fields that should not be empty.
func FieldsNotEmpty(fillStruct interface{}, fieldsName []string) error {
	if reflect.TypeOf(fillStruct).Kind() != reflect.Pointer {
		return ErrArgumentNotPointer{"fillStruct"}
	}
	of := reflect.ValueOf(fillStruct).Elem()
	typeOf := reflect.TypeOf(fillStruct).Elem()
	for i := 0; i < len(fieldsName); i++ {
		_, ok := typeOf.FieldByName(fieldsName[i])
		if ok {
			val := of.FieldByName(fieldsName[i])
			if val.IsNil() {
				return errors.New(fmt.Sprintf("field %s is empty", fieldsName[i]))
			}
		}
	}
	return nil
}
