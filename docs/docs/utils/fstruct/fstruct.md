## fstruct

#### CheckNotDefaultFields
Checks whether the values ​​of the structure fields are default.<br>
That is, if the field is not passed or initialized with standard values, for example, nil.
```golang
func CheckNotDefaultFields(objectPtr typeopr.IPtr) error {
	objectLink := objectPtr.Ptr()
	var rObjectValue reflect.Value
	var rObjectType reflect.Type
	if reflect.TypeOf(objectLink).Elem() == reflect.TypeOf(reflect.Value{}) {
		rObjectValue = objectLink.(*reflect.Value).Elem()
		rObjectType = objectLink.(*reflect.Value).Elem().Type()
	} else {
		rObjectValue = reflect.ValueOf(objectLink).Elem()
		rObjectType = reflect.TypeOf(objectLink).Elem()
	}
	for i := 0; i < rObjectType.NumField(); i++ {
		fieldType := rObjectType.Field(i)
		fieldValue := rObjectValue.Field(i)
		tag := fieldType.Tag.Get("notdef")
		if tag == "" {
			continue
		}
		reqiredValue, err := strconv.ParseBool(tag)
		if err != nil {
			return err
		}
		if reqiredValue {
			if fieldValue.IsZero() {
				return ErrStructFieldIsDefault{fieldType.Name}
			}
		}
	}
	return nil
}
```