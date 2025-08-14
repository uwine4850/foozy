## fmap

#### MergeMap
Merges two maps into one.<br>
For example, if you pass Map1 and Map2, Map2 data will be added to Map1.
```golang
func MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2) {
	for key, value := range map2 {
		(*map1)[key] = value
	}
}
```

#### MergeMapSync
Merges two maps into one.<br>
For example, if you pass Map1 and Map2, Map2 data will be added to Map1.<br>
Safe for use in asynchronous mode.
```golang
func MergeMapSync[T1 comparable, T2 any](mu *sync.Mutex, map1 *map[T1]T2, map2 map[T1]T2) {
	mu.Lock()
	defer mu.Unlock()
	for key, value := range map2 {
		(*map1)[key] = value
	}
}
```

#### YamlMapToStruct
Writes a yaml map to the structure.<br>
__IMPORTANT:__ the field of the structure to be written must have the 
yaml `tag:"<field_name>"`. This name must correspond to the name of the field in the targetMap structure.<br>
Works in depth, you can make as many attachments as you want.
```golang
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
```