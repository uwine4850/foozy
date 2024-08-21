package formmapper

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

type TypedMapper struct {
	Form                   *form.Form
	Output                 itypeopr.IPtr
	IgnoreNotSupConversion bool
	NilIfNotExist          []string
}

func (m *TypedMapper) Fill() error {
	of := FrmValueToOrderedForm(m.Form)
	if err := m.fillStruct(of); err != nil {
		return err
	}
	return nil
}

func (m *TypedMapper) fillStruct(OF *OrderedForm) error {
	output := m.Output.Ptr()
	outT := reflect.TypeOf(output).Elem()
	outV := reflect.ValueOf(output).Elem()
	for i := 0; i < outT.NumField(); i++ {
		fieldT := outT.Field(i)
		fieldV := outV.Field(i)
		tag := fieldT.Tag.Get("form")
		// Skip if the tag is not a form
		if tag == "" {
			continue
		}
		OFval, ok := OF.GetByName(tag)
		if !ok {
			if m.NilIfNotExist != nil && fslice.SliceContains(m.NilIfNotExist, tag) {
				continue
			} else {
				return form.ErrFormConvertFieldNotFound{Field: tag}
			}
		}
		if err := m.handleItem(&OFval, &fieldV, fieldT.Name, 0, fieldT.Tag.Get("empty")); err != nil {
			return err
		}
	}
	return nil
}

func (m *TypedMapper) handleItem(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, index int, emptyValue string) error {
	fieldT := fieldV.Type()
	switch fieldT.Kind() {
	case reflect.Int:
		if err := m.parseInt(val, fieldV, fieldName, index, emptyValue); err != nil {
			return err
		}
	case reflect.Float64:
		if err := m.parseFloat64(val, fieldV, fieldName, index, emptyValue); err != nil {
			return err
		}
	case reflect.String:
		if err := m.parseString(val, fieldV, fieldName, index, emptyValue); err != nil {
			return err
		}
	case reflect.Bool:
		if err := m.parseBool(val, fieldV, fieldName, index, emptyValue); err != nil {
			return err
		}
	case reflect.Struct:
		if fieldT == reflect.TypeOf(form.FormFile{}) {
			m.parseFile(val, fieldV, emptyValue)
		} else {
			return fmt.Errorf("conversion to %s type is not supported", fieldT.Kind().String())
		}
	case reflect.Slice:
		if err := m.parseSlice(val, fieldV, fieldT, fieldName, emptyValue); err != nil {
			return err
		}
	default:
		if !m.IgnoreNotSupConversion {
			return fmt.Errorf("conversion to %s type is not supported", fieldT.Kind().String())
		}
	}
	return nil
}

func (m *TypedMapper) parseInt(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, index int, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := m.fieldValue(stringSlice, fieldName, index, emptyValue)
	if err != nil {
		return err
	}
	if stringFieldValue != "" {
		intVal, err := strconv.Atoi(stringFieldValue)
		if err != nil {
			return err
		}
		fieldV.Set(reflect.ValueOf(intVal))
	}
	return nil
}

func (m *TypedMapper) parseFloat64(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, index int, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := m.fieldValue(stringSlice, fieldName, index, emptyValue)
	if err != nil {
		return err
	}
	if stringFieldValue != "" {
		floatVal, err := strconv.ParseFloat(stringFieldValue, 64)
		if err != nil {
			return err
		}
		fieldV.Set(reflect.ValueOf(floatVal))
	}
	return nil
}

func (m *TypedMapper) parseString(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, index int, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := m.fieldValue(stringSlice, fieldName, index, emptyValue)
	if err != nil {
		return err
	}
	fieldV.Set(reflect.ValueOf(stringFieldValue))
	return nil
}

func (m *TypedMapper) parseBool(val *[]OrderedFormValue, fieldV *reflect.Value, fieldName string, index int, emptyValue string) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	stringFieldValue, err := m.fieldValue(stringSlice, fieldName, index, emptyValue)
	if err != nil {
		return err
	}
	if stringFieldValue != "" {
		boolVal, err := strconv.ParseBool(stringFieldValue)
		if err != nil {
			return err
		}
		fieldV.Set(reflect.ValueOf(boolVal))
	}
	return nil
}

func (m *TypedMapper) parseFile(val *[]OrderedFormValue, fieldV *reflect.Value, emptyValue string) {
	if emptyValue != "" {
		panic("the field that accepts the file does not support the 'empty' tag")
	}
	fileSlice := reflect.ValueOf((*val)[0].Value).Interface().([]form.FormFile)
	if len(fileSlice) != 0 {
		fieldV.Set(reflect.ValueOf(fileSlice[0]))
	}
}

func (m *TypedMapper) parseSlice(val *[]OrderedFormValue, fieldV *reflect.Value, fieldT reflect.Type, fieldName string, emptyValue string) error {
	v := reflect.ValueOf((*val)[0].Value)
	if v.Type() == reflect.TypeOf([]form.FormFile{}) {
		fileSlice := v.Interface().([]form.FormFile)
		if len(fileSlice) != 0 {
			fieldV.Set(reflect.ValueOf(fileSlice))
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
			if err := m.handleItem(&[]OrderedFormValue{
				*newOF,
			}, &sliceItem, fieldName, i, emptyValue); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, sliceItem)
		}
		fieldV.Set(newSlice)
	}
	return nil
}

func (m *TypedMapper) fieldValue(stringSlice []string, fieldName string, index int, emptyValue string) (string, error) {
	var stringFieldValue string
	if len(stringSlice) != 0 && stringSlice[0] != "" {
		stringFieldValue = stringSlice[0]
	} else if emptyValue != "" {
		val, err := m.parseEmptyValue(fieldName, emptyValue, index)
		if err != nil {
			return "", err
		}
		stringFieldValue = val
	}
	return stringFieldValue, nil
}

func (m *TypedMapper) parseEmptyValue(fieldName string, emptyValue string, index int) (string, error) {
	switch emptyValue {
	case "-err":
		return "", ErrEmptyFieldIndex{Name: fieldName, Index: strconv.Itoa(index)}
	default:
		return emptyValue, nil
	}
}
