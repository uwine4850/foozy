package mapper

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/uwine4850/foozy/pkg/router/form"
)

const (
	EMPTY_VALUE = "v"
	EMPTY_ERROR = "-err"
)

const (
	NIL_SKIP = "-skip"
)

type FormConverter struct{}

func (fc FormConverter) handleItem(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, emptyValue string) error {
	fieldT := fieldV.Type()
	switch fieldT.Kind() {
	case reflect.Int:
		if err := fc.parseInt(val, fieldV, fieldName, emptyValue); err != nil {
			return err
		}
	case reflect.Float64:
		if err := fc.parseFloat64(val, fieldV, fieldName, emptyValue); err != nil {
			return err
		}
	case reflect.String:
		if err := fc.parseString(val, fieldV, fieldName, emptyValue); err != nil {
			return err
		}
	case reflect.Bool:
		if err := fc.parseBool(val, fieldV, fieldName, emptyValue); err != nil {
			return err
		}
	case reflect.Struct:
		if fieldT == formFileType {
			fc.parseFile(val, fieldV, emptyValue)
		} else {
			return &ErrUnsupportedStructField{FieldName: fieldName}
		}
	case reflect.Slice:
		if fieldT.Elem().Kind() == reflect.Struct && fieldT.Elem() != formFileType {
			return &ErrUnsupportedStructField{FieldName: fieldName}
		}
		if err := fc.parseSlice(val, fieldV, fieldT, fieldName, emptyValue); err != nil {
			return err
		}
	default:
		return &ErrUnsupportedStructField{FieldName: fieldName}
	}
	return nil
}

func (fc FormConverter) parseInt(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := fc.fieldStringValue(stringSlice, fieldName, emptyValue)
	if err != nil {
		return err
	}
	if stringFieldValue != "" {
		intVal, err := strconv.Atoi(stringFieldValue)
		if err != nil {
			return err
		}
		newV := reflect.ValueOf(intVal)
		if !fieldV.CanSet() || !newV.Type().AssignableTo(fieldV.Type()) {
			return &ErrFieldNotSettable{FieldName: fieldName}
		}
		fieldV.Set(newV)
	}
	return nil
}

func (fc FormConverter) parseFloat64(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := fc.fieldStringValue(stringSlice, fieldName, emptyValue)
	if err != nil {
		return err
	}
	if stringFieldValue != "" {
		floatVal, err := strconv.ParseFloat(stringFieldValue, 64)
		if err != nil {
			return err
		}
		newV := reflect.ValueOf(floatVal)
		if !fieldV.CanSet() || !newV.Type().AssignableTo(fieldV.Type()) {
			return &ErrFieldNotSettable{FieldName: fieldName}
		}
		fieldV.Set(newV)
	}
	return nil
}

func (fc FormConverter) parseString(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := fc.fieldStringValue(stringSlice, fieldName, emptyValue)
	if err != nil {
		return err
	}
	newV := reflect.ValueOf(stringFieldValue)
	if !fieldV.CanSet() || !newV.Type().AssignableTo(fieldV.Type()) {
		return &ErrFieldNotSettable{FieldName: fieldName}
	}
	fieldV.Set(newV)
	return nil
}

func (fc FormConverter) parseBool(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := fc.fieldStringValue(stringSlice, fieldName, emptyValue)
	if err != nil {
		return err
	}
	if stringFieldValue != "" {
		boolVal, err := strconv.ParseBool(stringFieldValue)
		if err != nil {
			return err
		}
		newV := reflect.ValueOf(boolVal)
		if !fieldV.CanSet() || !newV.Type().AssignableTo(fieldV.Type()) {
			return &ErrFieldNotSettable{FieldName: fieldName}
		}
		fieldV.Set(newV)
	}
	return nil
}

func (fc FormConverter) parseFile(val *[]OrderedFormValue, fieldV *reflect.Value, emptyValue string) {
	if emptyValue != "" {
		panic("the field that accepts the file does not support the 'empty' tag")
	}
	fileSlice := reflect.ValueOf((*val)[0].Value).Interface().([]form.FormFile)
	if len(fileSlice) != 0 {
		fieldV.Set(reflect.ValueOf(fileSlice[0]))
	}
}

func (fc FormConverter) parseSlice(val *[]OrderedFormValue, fieldV *reflect.Value, fieldT reflect.Type, fieldName string, emptyValue string) error {
	v := reflect.ValueOf((*val)[0].Value)
	if v.Type() == reflect.TypeOf([]form.FormFile{}) {
		fileSlice := v.Interface().([]form.FormFile)
		if len(fileSlice) != 0 {
			newV := reflect.ValueOf(fileSlice)
			if !fieldV.CanSet() || !newV.Type().AssignableTo(fieldV.Type()) {
				return &ErrFieldNotSettable{FieldName: fieldName}
			}
			fieldV.Set(newV)
		}
		return nil
	}
	stringSlice := v.Interface().([]string)
	if len(stringSlice) != 0 {
		sliceType := fieldT.Elem()
		newSlice := reflect.MakeSlice(fieldT, 0, len(stringSlice))
		for i := 0; i < len(stringSlice); i++ {
			newOF := &OrderedFormValue{
				Name:  "",
				Value: []string{stringSlice[i]},
			}
			sliceItem := reflect.New(sliceType).Elem()
			if err := fc.handleItem(&[]OrderedFormValue{
				*newOF,
			}, &sliceItem, fieldName, emptyValue); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, sliceItem)
		}
		if !fieldV.CanSet() || !newSlice.Type().AssignableTo(fieldV.Type()) {
			return &ErrFieldNotSettable{FieldName: fieldName}
		}
		fieldV.Set(newSlice)
	}
	return nil
}

func (fc FormConverter) fieldStringValue(stringSlice []string, fieldName string, emptyValue string) (string, error) {
	if len(stringSlice) != 0 && stringSlice[0] != "" {
		return stringSlice[0], nil
	} else if emptyValue != "" {
		status, val := fc.parseEmptyArgs(emptyValue)
		switch status {
		case EMPTY_ERROR:
			return "", &ErrEmptyFieldIndex{Name: fieldName, Index: "undefined"}
		case EMPTY_VALUE:
			return val, nil
		}
	}
	return "", nil
}

func (fc FormConverter) parseEmptyArgs(emptyValue string) (status string, val string) {
	switch emptyValue {
	case EMPTY_ERROR:
		return EMPTY_ERROR, ""
	default:
		return EMPTY_VALUE, emptyValue
	}
}

type ErrUnsupportedStructField struct {
	FieldName string
}

func (e ErrUnsupportedStructField) Error() string {
	return fmt.Sprintf("%s field conversion is not possible due to unsupported type", e.FieldName)
}

type ErrFieldNotSettable struct {
	FieldName string
}

func (e ErrFieldNotSettable) Error() string {
	return fmt.Sprintf("field %s is not settable", e.FieldName)
}

type ErrEmptyFieldIndex struct {
	Name  string
	Index string
}

func (e ErrEmptyFieldIndex) Error() string {
	return fmt.Sprintf("the %s field value at index %s is empty", e.Name, e.Index)
}
