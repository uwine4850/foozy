## typeopr
This package contains functions for working with data types.

#### IsPointer
Checks if the value is a pointer.
```golang
func IsPointer(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Pointer
}
```

#### PtrIsStruct
Checks if the value is a pointer to a structure.
```golang
func PtrIsStruct(a any) bool {
	return reflect.TypeOf(a).Elem().Kind() == reflect.Struct
}
```

#### IsEmpty
Checks if the value is empty.
```golang
func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		return v.NumField() == 0
	case reflect.Invalid:
		return true
	}
	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}
```

#### AnyToBytes
Converts input data into bytes.
```golang
func AnyToBytes(value interface{}) ([]byte, error) {
	var buf bytes.Buffer

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("%v", val.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(fmt.Sprintf("%v", val.Uint()))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("%v", val.Float()))
	case reflect.String:
		buf.WriteString(val.String())
	default:
		return nil, fmt.Errorf("unsupported convert type %s", val.Kind().String())
	}
	return buf.Bytes(), nil
}
```

### Ptr
An object designed to store a pointer. This object guarantees that it will store a pointer.

#### Ptr.New
Creating a new object `Ptr`.<br>
It is also guaranteed that a pointer will be stored inside this object.
```golang
func (p Ptr) New(value interface{}) IPtr {
	if !IsPointer(value) {
		panic(ErrValueNotPointer{Value: fmt.Sprintf("Ptr<%s>", reflect.TypeOf(value))})
	}
	p.value = value
	return p
}
```

#### Prt.Ptr
Returns a pointer that is stored in the object.
```golang
func (p Ptr) Ptr() interface{} {
	return p.value
}
```

#### IsImplementInterface
Determines whether an object uses an interface.<br>
If `reflect.Value` is used, you can use direct passing or passing by pointer, that is, 
passing a pointer to a pointer. How to transmit this data depends on the situation.<br>
```golang
func IsImplementInterface(objectPtr IPtr, interfaceType interface{}) bool {
	object := objectPtr.Ptr()
	// If the type of data passed directly is the desired interface.
	if reflect.TypeOf(object) == reflect.TypeOf(interfaceType) {
		return true
	}
	var objType reflect.Type
	if reflect.TypeOf(object).Elem() == reflect.TypeOf(reflect.Value{}) {
		objType = object.(*reflect.Value).Type()
	} else {
		objType = reflect.TypeOf(object)
	}
	intrfcType := reflect.TypeOf(interfaceType).Elem()
	return objType.Implements(intrfcType)
}
```
Usage example:
```golang
object := MyObject{}
IsImplementInterface(typeopr.Ptr{}.New(&object), (*MyInterface)(nil))
```

#### GetReflectValue
Get `reflect.Value` from the passed value.
```golang
func GetReflectValue[T any](target *T) reflect.Value {
	var v reflect.Value
	if reflect.TypeOf(*target) == typeReflectValue {
		v = *reflect.ValueOf(target).Interface().(*reflect.Value)
	} else {
		v = reflect.ValueOf(target).Elem()
	}
	return v
}
```