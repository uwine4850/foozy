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
		if err := m.handleItem(&OFval, &fieldV); err != nil {
			return err
		}
	}
	return nil
}

func (m *TypedMapper) handleItem(val *[]OrderedFormValue, fieldV *reflect.Value) error {
	fieldT := fieldV.Type()
	switch fieldT.Kind() {
	case reflect.Int:
		if err := m.parseInt(val, fieldV); err != nil {
			return err
		}
	case reflect.Float64:
		if err := m.parseFloat64(val, fieldV); err != nil {
			return err
		}
	case reflect.String:
		m.parseString(val, fieldV)
	case reflect.Bool:
		if err := m.parseBool(val, fieldV); err != nil {
			return err
		}
	case reflect.Struct:
		if fieldT == reflect.TypeOf(form.FormFile{}) {
			m.parseFile(val, fieldV)
		} else {
			return fmt.Errorf("conversion to %s type is not supported", fieldT.Kind().String())
		}
	case reflect.Slice:
		if err := m.parseSlice(val, fieldV, fieldT); err != nil {
			return err
		}
	default:
		if !m.IgnoreNotSupConversion {
			return fmt.Errorf("conversion to %s type is not supported", fieldT.Kind().String())
		}
	}
	return nil
}

func (m *TypedMapper) parseInt(val *[]OrderedFormValue, fieldV *reflect.Value) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	if len(stringSlice) != 0 {
		intVal, err := strconv.Atoi(stringSlice[0])
		if err != nil {
			return err
		}
		fieldV.Set(reflect.ValueOf(intVal))
	}
	return nil
}

func (m *TypedMapper) parseFloat64(val *[]OrderedFormValue, fieldV *reflect.Value) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	if len(stringSlice) != 0 {
		floatVal, err := strconv.ParseFloat(stringSlice[0], 64)
		if err != nil {
			return err
		}
		fieldV.Set(reflect.ValueOf(floatVal))
	}
	return nil
}

func (m *TypedMapper) parseString(val *[]OrderedFormValue, fieldV *reflect.Value) {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	if len(stringSlice) != 0 {
		fieldV.Set(reflect.ValueOf(stringSlice[0]))
	}
}

func (m *TypedMapper) parseBool(val *[]OrderedFormValue, fieldV *reflect.Value) error {
	stringSlice := reflect.ValueOf((*val)[0].Value).Interface().([]string)
	if len(stringSlice) != 0 {
		boolVal, err := strconv.ParseBool(stringSlice[0])
		if err != nil {
			return err
		}
		fieldV.Set(reflect.ValueOf(boolVal))
	}
	return nil
}

func (m *TypedMapper) parseFile(val *[]OrderedFormValue, fieldV *reflect.Value) {
	fileSlice := reflect.ValueOf((*val)[0].Value).Interface().([]form.FormFile)
	if len(fileSlice) != 0 {
		fieldV.Set(reflect.ValueOf(fileSlice[0]))
	}
}

func (m *TypedMapper) parseSlice(val *[]OrderedFormValue, fieldV *reflect.Value, fieldT reflect.Type) error {
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
			}, &sliceItem); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, sliceItem)
		}
		fieldV.Set(newSlice)
	}
	return nil
}
