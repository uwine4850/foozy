package form

import (
	"errors"
	"fmt"
	"github.com/uwine4850/foozy/pkg/ferrors"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"mime/multipart"
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
func FillStructFromForm(frm interfaces.IForm, fill interface{}) error {
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
		err := checkFieldType(field)
		if err != nil {
			return err
		}
		value := v.Field(i)
		tag := field.Tag.Get("form")
		// Skip if the tag is not a form
		if tag == "" {
			continue
		}
		formMapItem, ok := formMap[tag]
		if !ok {
			return ErrFormConvertFieldNotFound{tag}
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

// FrmValueToMap Converts the form to a map.
func FrmValueToMap(frm interfaces.IForm) map[string]interface{} {
	formMap := make(map[string]interface{})
	multipartForm := frm.GetMultipartForm()
	if multipartForm != nil {
		for name, value := range multipartForm.Value {
			formMap[name] = value
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
