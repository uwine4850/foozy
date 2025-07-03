package fmap

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

// MergeMap merges two maps into one.
// For example, if you pass Map1 and Map2, Map2 data will be added to Map1.
func MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2) {
	for key, value := range map2 {
		(*map1)[key] = value
	}
}

func MergeMapSync[T1 comparable, T2 any](mu *sync.Mutex, map1 *map[T1]T2, map2 map[T1]T2) {
	mu.Lock()
	defer mu.Unlock()
	for key, value := range map2 {
		(*map1)[key] = value
	}
}

// Compare map values. It is important that the keys and values ​​match.
// exclude - keys that do not need to be taken into account.
// func Compare[T1 comparable, T2 comparable](map1 *map[T1]T2, map2 *map[T1]T2, exclude []T1) bool {
// 	for key, value := range *map1 {
// 		if exclude != nil && fslice.SliceContains(exclude, key) {
// 			continue
// 		}
// 		value2, ok := (*map2)[key]
// 		if !ok {
// 			return false
// 		} else {
// 			if value != value2 {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// YamlMapToStruct writes a yaml map to the structure.
// IMPOrTANT: the field of the structure to be written must have the
// yaml tag:"<field_name>". This name must correspond to the name of the
// field in the targetMap structure.
// Works in depth, you can make as many attachments as you want.
func YamlMapToStruct(targetMap *map[string]interface{}, targetStruct typeopr.IPtr) error {
	for mFieldName, mFieldValue := range *targetMap {
		var sValue reflect.Value
		var sType reflect.Type
		if reflect.DeepEqual(reflect.TypeOf(targetStruct.Ptr()).Elem(), reflect.TypeOf(reflect.Value{})) {
			sValue = *targetStruct.Ptr().(*reflect.Value)
			sType = sValue.Type()
		} else {
			sValue = reflect.ValueOf(targetStruct.Ptr()).Elem()
			sType = reflect.TypeOf(targetStruct.Ptr()).Elem()
		}

		for i := 0; i < sType.NumField(); i++ {
			if sType.Field(i).Tag.Get("yaml") == mFieldName {
				if sValue.CanSet() {
					fieldValue := sValue.Field(i)
					mapFieldValue := reflect.ValueOf(mFieldValue)
					if err := convertYamlField(&fieldValue, &mapFieldValue); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("the %s field cannot be set to a value", sType.Field(i).Name)
				}
			}
		}
	}
	return nil
}

func convertYamlField(fieldValue *reflect.Value, yamlValue *reflect.Value) error {
	switch fieldValue.Kind() {
	case reflect.Struct:
		nextTargetMap := yamlValue.Interface().(map[string]interface{})
		if err := YamlMapToStruct(&nextTargetMap, typeopr.Ptr{}.New(fieldValue)); err != nil {
			return err
		}
	case reflect.Slice:
		yamlSlice := yamlValue.Interface().([]interface{})
		newSlice := reflect.MakeSlice(fieldValue.Type(), len(yamlSlice), len(yamlSlice))
		for i := 0; i < len(yamlSlice); i++ {
			if reflect.TypeOf(yamlSlice[i]).Kind() == reflect.Slice {
				yamlSliceElemValue := reflect.ValueOf(yamlSlice[i])
				newInnerSlice := reflect.New(fieldValue.Type().Elem()).Elem()
				if err := convertYamlField(&newInnerSlice, &yamlSliceElemValue); err != nil {
					return err
				}
				newSlice.Index(i).Set(newInnerSlice)
			} else {
				field := newSlice.Index(i)
				val := reflect.ValueOf(yamlSlice[i])
				if err := convertYamlField(&field, &val); err != nil {
					return err
				}
			}
		}
		fieldValue.Set(newSlice)
	case reflect.Map:
		return errors.New("map data type is not supported, you need to use a structure")
	default:
		if fieldValue.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(yamlValue.Convert(fieldValue.Type()))
		} else {
			return fmt.Errorf("cannot assign value of type %s to field of type %s", yamlValue.String(), fieldValue.Type().String())
		}
	}
	return nil
}
