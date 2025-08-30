## dbmapper

#### FillStructSliceFromDb
Fills a slice with data from the database.
It uses `FillStructFromDb` function for filling.
```golang
func FillStructSliceFromDb[T any](slice *[]T, dbRes *[]map[string]interface{}) error {
	if len(*slice) != len(*dbRes) {
		return errors.New("the length of the fill slice is not the same as the length of the data slice")
	}
	for i := 0; i < len(*dbRes); i++ {
		if err := FillStructFromDb(&(*slice)[i], &(*dbRes)[i]); err != nil {
			return err
		}
	}
	return nil
}
```

#### FillStructFromDb
Fills the structure with data from the database.<br>
It needs the `db:"<field_name>"` tag to work properly. The name of the 
tag must match the name of the column.
If there is no tag, the field is skipped.

Caches the structure using the implemented [RawObject](/mapper/mapper/#rawobject) interface.
This means that all subsequent accesses to this structure will be faster.
```golang
func FillStructFromDb[T any](fillStruct *T, dbRes *map[string]interface{}) error {
	v := typeopr.GetReflectValue(fillStruct)
	raw := LoadSomeRawObjectFromCache(v, &dbRawCache, namelib.TAGS.DB_MAPPER_NAME)

	for name, f := range *raw.Fields() {
		field := v.FieldByName(f.Name)
		data, ok := (*dbRes)[name]
		if ok {
			// Processing DB_MAPPER_EMPTY tag.
			if data == nil {
				emptyVal := f.Tag.Get(namelib.TAGS.DB_MAPPER_EMPTY)
				if emptyVal != "" {
					if emptyVal == "-error" {
						return typeopr.ErrValueIsEmpty{Value: name}
					}
					newByteData, err := DC.dbValueConversionToByte(emptyVal)
					if err != nil {
						return err
					}
					newData := reflect.ValueOf(newByteData).Interface()
					if err := DC.convertDBType(&field, &f.Tag, &newData); err != nil {
						return err
					}
				}
			} else {
				if err := DC.convertDBType(&field, &f.Tag, &data); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
```

#### ParamsValueFromDbStruct
Creates a map from a structure that describes the table.
To work correctly, you need a completed structure, and the required fields must have the `db:"<column name>"` tag.
```golang
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
```