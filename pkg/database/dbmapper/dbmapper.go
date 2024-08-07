package dbmapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

type Mapper struct {
	DatabaseResult []map[string]interface{}
	Output         interface{}
}

func NewMapper(databaseResult []map[string]interface{}, output interface{}) Mapper {
	return Mapper{DatabaseResult: databaseResult, Output: output}
}

func (m *Mapper) SetOutput(output interface{}) {
	m.Output = output
}

func (m *Mapper) Fill() error {
	outType, err := m.outputType()
	if err != nil {
		return err
	}
	sliceType := reflect.TypeOf(m.Output).Elem().Elem()
	newOutputSlice := reflect.MakeSlice(reflect.TypeOf(m.Output).Elem(), 0, len(m.DatabaseResult))
	switch outType {
	case reflect.Struct:
		if err := m.fillStruct(sliceType, &newOutputSlice); err != nil {
			return err
		}
		reflect.ValueOf(m.Output).Elem().Set(newOutputSlice)
	case reflect.Map:
		if err := m.fillMap(sliceType, &newOutputSlice); err != nil {
			return err
		}
		reflect.ValueOf(m.Output).Elem().Set(newOutputSlice)
	default:
		return fmt.Errorf("mapping for %s type is not supported", sliceType)
	}
	return nil
}

func (m *Mapper) outputType() (reflect.Kind, error) {
	if !typeopr.IsPointer(m.Output) {
		return reflect.Invalid, typeopr.ErrValueNotPointer{Value: "Output"}
	}
	typeOf := reflect.TypeOf(m.Output).Elem()
	if typeOf.Kind() != reflect.Slice {
		return reflect.Invalid, errors.New("field Output must be a slice")
	} else {
		return typeOf.Elem().Kind(), nil
	}
}

func (m *Mapper) fillStruct(typeOut reflect.Type, newOutputSlice *reflect.Value) error {
	for i := 0; i < len(m.DatabaseResult); i++ {
		fill := reflect.New(typeOut).Elem()
		if err := FillReflectValueFromDb(m.DatabaseResult[i], &fill); err != nil {
			return err
		}
		*newOutputSlice = reflect.Append(*newOutputSlice, fill)
	}
	return nil
}

func (m *Mapper) fillMap(typeOut reflect.Type, newOutputSlice *reflect.Value) error {
	for i := 0; i < len(m.DatabaseResult); i++ {
		fill := reflect.MakeMap(typeOut)
		f, ok := fill.Interface().(map[string]string)
		if !ok {
			return errors.New("cannot assert the Output type as a map[string]string type")
		}
		if err := FillMapFromDb(m.DatabaseResult[i], &f); err != nil {
			return err
		}
		*newOutputSlice = reflect.Append(*newOutputSlice, fill)
	}
	return nil
}
