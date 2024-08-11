package formmapper

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

// FillStructFromForm A method that fills the structure with data from the form.
// The structure should always be passed as a pointer.
// For correct work it is necessary to specify "form" tag for each field of the structure. For example, `form:<form field name>`.
// Structure fields can be of two types only:
// []FormFile - form files.
// []string - all other data.
func FillStructFromForm(frm *form.Form, fillPtr itypeopr.IPtr, nilIfNotExist []string) error {
	fillStruct := fillPtr.Ptr()
	if !typeopr.PtrIsStruct(fillStruct) {
		return typeopr.ErrParameterNotStruct{Param: "fillStruct"}
	}
	orderedForm := FrmValueToOrderedForm(frm)
	t := reflect.TypeOf(fillStruct).Elem()
	v := reflect.ValueOf(fillStruct).Elem()
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
		orderedFormValue, ok := orderedForm.GetByName(tag)
		if !ok {
			if nilIfNotExist != nil && fslice.SliceContains(nilIfNotExist, tag) {
				continue
			} else {
				return form.ErrFormConvertFieldNotFound{Field: tag}
			}
		}
		for i := 0; i < len(orderedFormValue); i++ {
			//Be sure to first check the setting of the string value.
			// This happens because a field without a file loaded is considered a string.
			// Therefore, an empty []FormFile field will be of type []string.
			// If the field is suitable for a string, then you need to skip the iteration to avoid false positives when processing files.
			ok, err = setFormString(field, orderedFormValue[i].Value, value, field.Tag.Get("empty"))
			if err != nil {
				return err
			}
			if ok {
				continue
			}
			if err := setFormFile(field, orderedFormValue[i].Value, value, field.Tag.Get("empty")); err != nil {
				return err
			}
		}
	}
	return nil
}

// FillReflectValueFromForm fills the structure with data from the form.
// It works and does everything the same as the FillStructFromForm function.
// The only difference is that this function accepts data in *reflect.Value format.
func FillReflectValueFromForm(frm *form.Form, fillValue *reflect.Value, nilIfNotExist []string) error {
	orderedForm := FrmValueToOrderedForm(frm)
	fillType := fillValue.Type()
	for i := 0; i < fillType.NumField(); i++ {
		field := fillType.Field(i)
		value := fillValue.Field(i)

		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}
		orderedFormValue, ok := orderedForm.GetByName(tag)
		if !ok {
			if nilIfNotExist != nil && fslice.SliceContains(nilIfNotExist, tag) {
				continue
			} else {
				return form.ErrFormConvertFieldNotFound{Field: tag}
			}
		}
		for i := 0; i < len(orderedFormValue); i++ {
			//Be sure to first check the setting of the string value.
			// This happens because a field without a file loaded is considered a string.
			// Therefore, an empty []FormFile field will be of type []string.
			// If the field is suitable for a string, then you need to skip the iteration to avoid false positives when processing files.
			ok, err := setFormString(field, orderedFormValue[i].Value, value, field.Tag.Get("empty"))
			if err != nil {
				return err
			}
			if ok {
				continue
			}
			if err := setFormFile(field, orderedFormValue[i].Value, value, field.Tag.Get("empty")); err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckExtension Check if the file resolution matches the expected one. Can only be used with a structure already
// filled out in the form.
// To work, you need to add an ext tag with the necessary extensions (if there are many, separated by commas).
// For example, ext:".jpg .jpeg .png".
func CheckExtension(fillPtr itypeopr.IPtr) error {
	fillStruct := fillPtr.Ptr()
	if !typeopr.PtrIsStruct(fillStruct) {
		return typeopr.ErrParameterNotStruct{Param: "fillStruct"}
	}

	fillType, fillValue := getStructElemFromFillableStruct(fillStruct)

	for i := 0; i < fillType.NumField(); i++ {
		tag := fillType.Field(i).Tag.Get("ext")
		if tag == "" {
			continue
		}
		field := fillType.Field(i)
		if field.Type.Elem() != reflect.TypeOf(form.FormFile{}) {
			panic("the ext tag can only be added to fields whose type is form.FormFile")
		}
		extension := strings.Split(tag, " ")
		files := fillValue.Field(i).Interface().([]form.FormFile)
		for i := 0; i < len(files); i++ {
			if !fslice.SliceContains(extension, filepath.Ext(files[i].Header.Filename)) {
				return ErrExtensionNotMatch{Field: field.Name}
			}
		}
	}
	return nil
}

// FieldsName returns a list of field names of the filled structure.
func FieldsName(fillPtr itypeopr.IPtr, exclude []string) ([]string, error) {
	fillStruct := fillPtr.Ptr()
	if !typeopr.IsPointer(fillStruct) {
		return nil, typeopr.ErrValueNotPointer{Value: "fillStruct"}
	}
	if !typeopr.PtrIsStruct(fillStruct) {
		return nil, typeopr.ErrParameterNotStruct{Param: "fillStruct"}
	}

	fillType, _ := getStructElemFromFillableStruct(fillStruct)
	var names []string
	for i := 0; i < fillType.NumField(); i++ {
		field := fillType.Field(i)
		if exclude != nil && fslice.SliceContains(exclude, field.Name) {
			continue
		}
		names = append(names, field.Name)
	}
	return names, nil
}

// FieldsNotEmpty checks the specified fields of the structure for emptiness.
// fieldsName - slice with exact names of STRUCTURE fields that should not be empty.
// Optimized to work even if the FillableFormStruct contains a structure with type *reflect.Value.
func FieldsNotEmpty(fillPtr itypeopr.IPtr, fieldsName []string) error {
	fillStruct := fillPtr.Ptr()
	if !typeopr.IsPointer(fillStruct) {
		return typeopr.ErrValueNotPointer{Value: "fillStruct"}
	}
	if !typeopr.PtrIsStruct(fillStruct) {
		return typeopr.ErrParameterNotStruct{Param: "fillStruct"}
	}

	fillType, fillValue := getStructElemFromFillableStruct(fillStruct)
	for i := 0; i < len(fieldsName); i++ {
		_, ok := fillType.FieldByName(fieldsName[i])
		if ok {
			val := fillValue.FieldByName(fieldsName[i])
			if val.IsNil() {
				return fmt.Errorf("field %s is empty", fieldsName[i])
			}
			if val.Len() != 0 && val.Type().Elem().Kind() == reflect.String && fslice.AllStringItemsEmpty(val.Interface().([]string)) {
				return fmt.Errorf("field %s is empty", fieldsName[i])
			}
		}
	}
	return nil
}

func getStructElemFromFillableStruct(fillStruct interface{}) (reflect.Type, reflect.Value) {
	var fillValue reflect.Value
	var fillType reflect.Type
	// If fillStruct is reflect.Value.
	if reflect.TypeOf(fillStruct) == reflect.TypeOf(&reflect.Value{}) {
		fv := fillStruct.(*reflect.Value)
		fillValue = reflect.Indirect(*fv)
		fillType = fillValue.Type()
	} else {
		// If fillStruct is default struct.
		fillValue = reflect.ValueOf(fillStruct).Elem()
		fillType = reflect.TypeOf(fillStruct).Elem()
	}
	return fillType, fillValue
}

// checkFieldType checks if the type of field to be filled is correct.
func checkFieldType(field reflect.StructField) error {
	if field.Type.Kind() != reflect.Slice {
		return fmt.Errorf("the %s field should be slice", field.Name)
	}
	if field.Type.Elem() != reflect.TypeOf("") && field.Type.Elem() != reflect.TypeOf(form.FormFile{}) {
		return fmt.Errorf("the %s field can only be of two types: []string or []FormFile", field.Name)
	}
	return nil
}

func setFormFile(field reflect.StructField, formValue interface{}, value reflect.Value, emptyTag string) error {
	if reflect.DeepEqual(field.Type, reflect.TypeOf([]form.FormFile{})) && reflect.TypeOf(formValue) == reflect.TypeOf([]form.FormFile{}) {
		formType, ok := formValue.([]form.FormFile)
		if !ok {
			return typeopr.ErrConvertType{Type1: reflect.TypeOf(formValue).String(), Type2: "[]FormFile"}
		}
		value.Set(reflect.ValueOf(formType))
	} else {
		if emptyTag != "" {
			for i := 0; i < reflect.ValueOf(formValue).Len(); i++ {
				if reflect.ValueOf(formValue).Index(i).IsZero() {
					if err := emptyOperationsFile(emptyTag, field.Name); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func setFormString(field reflect.StructField, formValue interface{}, value reflect.Value, emptyTag string) (bool, error) {
	isStringValue := false
	if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
		isStringValue = true
		formType, ok := formValue.([]string)
		if !ok {
			return isStringValue, typeopr.ErrConvertType{Type1: reflect.TypeOf(formValue).String(), Type2: "string"}
		}
		if emptyTag != "" {
			for i := 0; i < len(formType); i++ {
				if formType[i] == "" {
					if err := emptyOperations(emptyTag, field.Name, &formType[i], i); err != nil {
						return isStringValue, err
					}
				}
			}
		}
		value.Set(reflect.ValueOf(formType))
	}
	return isStringValue, nil
}

func emptyOperations(emptyTagValue string, fieldName string, value *string, index int) error {
	if emptyTagValue == "" {
		return nil
	}
	switch emptyTagValue {
	case "-err":
		return ErrEmptyFieldIndex{Name: fieldName, Index: strconv.Itoa(index)}
	default:
		*value = emptyTagValue
	}
	return nil
}

func emptyOperationsFile(emptyTagValue string, fieldName string) error {
	if emptyTagValue == "" {
		return nil
	}
	switch emptyTagValue {
	case "-err":
		return ErrEmptyFieldIndex{Name: fieldName, Index: "unkown"}
	default:
		return fmt.Errorf("the FormFile field does not support a default value")
	}
}

type ErrExtensionNotMatch struct {
	Field string
}

func (e ErrExtensionNotMatch) Error() string {
	return fmt.Sprintf("The extension of the %s field does not match what is expected.", e.Field)
}

type ErrEmptyFieldIndex struct {
	Name  string
	Index string
}

func (e ErrEmptyFieldIndex) Error() string {
	return fmt.Sprintf("the %s field value at index %s is empty", e.Name, e.Index)
}
