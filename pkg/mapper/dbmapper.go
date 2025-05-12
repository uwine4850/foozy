package mapper

import (
	"errors"
	"reflect"

	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

type DbRawStruct struct {
	_type  reflect.Type
	fields *map[string]reflect.StructField
}

func (s *DbRawStruct) Type() reflect.Type {
	return s._type
}

func (s *DbRawStruct) Fields() *map[string]reflect.StructField {
	return s.fields
}

var typeReflectValue = reflect.TypeOf(&reflect.Value{})

func NewDBRawStruct[T any](target *T) RawStruct {
	var t reflect.Type
	if reflect.DeepEqual(reflect.TypeOf(target), typeReflectValue) {
		v := *reflect.ValueOf(target).Interface().(*reflect.Value)
		t = v.Type()
	} else {
		t = reflect.TypeOf(target).Elem()
	}
	dbRawStruct := &DbRawStruct{}
	fields := make(map[string]reflect.StructField)
	dbRawStruct._type = t
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldDbTagName := field.Tag.Get(namelib.TAGS.DB_MAPPER_NAME)
		// Skip field if tag empty.
		if fieldDbTagName == "" {
			continue
		}
		fields[fieldDbTagName] = field
	}
	dbRawStruct.fields = &fields
	return dbRawStruct
}

func FillStructSliceFromDb[T any](rawStruct RawStruct, slice *[]T, dbRes *[]map[string]interface{}) error {
	if len(*slice) != len(*dbRes) {
		return errors.New("the length of the fill slice is not the same as the length of the data slice")
	}
	for i := 0; i < len(*dbRes); i++ {
		if err := FillStructFromDb(rawStruct, typeopr.Ptr{}.New(&(*slice)[i]), &(*dbRes)[i]); err != nil {
			return err
		}
	}
	return nil
}

func FillStructFromDb(rawStruct RawStruct, fillStructPtr typeopr.IPtr, dbRes *map[string]interface{}) error {
	fillStruct := fillStructPtr.Ptr()
	var fillStructValue reflect.Value
	if reflect.DeepEqual(reflect.TypeOf(fillStruct), typeReflectValue) {
		fillStructValue = *fillStruct.(*reflect.Value)
	} else {
		fillStructValue = reflect.ValueOf(fillStruct).Elem()
	}

	for name, f := range *rawStruct.Fields() {
		field := fillStructValue.FieldByName(f.Name)
		data, ok := (*dbRes)[name]
		if ok {
			// Processing DB_MAPPER_EMPTY tag.
			if data == nil {
				emptyVal := f.Tag.Get(namelib.TAGS.DB_MAPPER_EMPTY)
				if emptyVal != "" {
					if emptyVal == "-error" {
						return typeopr.ErrValueIsEmpty{Value: name}
					}
					newByteData, err := dbValueConversionToByte(emptyVal)
					if err != nil {
						return err
					}
					newData := reflect.ValueOf(newByteData).Interface()
					if err := convertDBType(&field, &f.Tag, &newData); err != nil {
						return err
					}
				}
			} else {
				if err := convertDBType(&field, &f.Tag, &data); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ParamsValueFromDbStruct creates a map from a structure that describes the table.
// To work correctly, you need a completed structure, and the required fields must have the `db:"<column name>"` tag.
func ParamsValueFromDbStruct(filledStructurePtr typeopr.IPtr, nilIfEmpty []string) (map[string]any, error) {
	structure := filledStructurePtr.Ptr()
	if !typeopr.PtrIsStruct(structure) {
		return nil, typeopr.ErrParameterNotStruct{Param: "structure"}
	}
	outputParamsMap := make(map[string]any)

	typeof := reflect.TypeOf(structure).Elem()
	valueof := reflect.ValueOf(structure).Elem()
	for i := 0; i < typeof.NumField(); i++ {
		fieldValue := valueof.Field(i)
		dbColName := typeof.Field(i).Tag.Get(namelib.TAGS.DB_MAPPER_NAME)
		if dbColName == "" {
			continue
		}
		if fslice.SliceContains(nilIfEmpty, dbColName) && fieldValue.IsZero() {
			outputParamsMap[dbColName] = nil
		} else {
			outputParamsMap[dbColName] = fieldValue.Interface()
		}
	}
	return outputParamsMap, nil
}
