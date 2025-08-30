## mapper
This package is used for mapping data to a data structure.<br>
The reflect package is actively used for the mapper's operation. Therefore, the key feature of the package is the use of caching.

The `LoadSomeRawObjectFromCache` function is used to create and subsequently load the cache. Example of use:
```golang
var rawCache sync.Map
...
v := typeopr.GetReflectValue(out)
raw := LoadSomeRawObjectFromCache(v, &rawCache, namelib.TAGS.FORM_MAPPER_NAME)
...
```

There are several types of mappers:

[dbmapper](/mapper/dbmapper)<br>
[formmapper](/mapper/formmapper)<br>
[dtomapper](/mapper/dtomapper)<br>

### RawObject
RawObject is used to store data of type undefined structure.

This structure is intended to optimize the work of the reflect package with the structure.<br>
The type of the structure itself and its field types are stored here. This object reduces
the load on the structure parser, now it is not necessary to loop through all fields.
Also this object is convenient because it can store only selected fields of the structure,
which makes it more convenient to use.

Fields are passed by the map so that it is possible to give them a custom tag name. If the tag
is not needed, the key can be the field name.

__IMPORTANT__: to get better performance you should store an instance of this object in a separate immutable variable.
```golang
type RawObject interface {
	// Type returns the type of the stored object.
	Type() reflect.Type
	// Fields returns the selected fields of the stored object.
	// It can store not all fields, but only those added by the implementation.
	Fields() *map[string]reflect.StructField
}
```

### SomeRawObject
Object stores structure data to be filled with data from the data base.
It implements the RawStruct interface.
```golang
type SomeRawObject struct {
	typ    reflect.Type
	fields *map[string]reflect.StructField
}
```

#### NewSomeRawObjectWithTag
Creates and fills a new instance of `SomeRawObject` from a given object.
Accepts an object as a direct instance or reflect.Value object.

Only fields that have the `<tag_name>:<field_name>` tag will be stored.
This tag must contain the names of the column in the table, for which
the structure field is intended. The names must exactly match.
If there is no tag, the field will be simply skipped.
```golang
func NewSomeRawObjectWithTag[T any](target *T, tagName string) RawObject {
	t := typeopr.GetReflectValue(target).Type()
	fields := make(map[string]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldDbTagName := field.Tag.Get(tagName)
		// Skip field if tag empty.
		if fieldDbTagName == "" {
			continue
		}
		fields[fieldDbTagName] = field
	}
	return &SomeRawObject{
		typ:    t,
		fields: &fields,
	}
}
```

#### LoadSomeRawObjectFromCache
Loads an object from the selected cache.
If the object is not in the cache, creates a RawObject for it and sets it.
```golang
func LoadSomeRawObjectFromCache(objectValue reflect.Value, rawCache *sync.Map, tagName string) RawObject {
	var raw RawObject
	objectType := objectValue.Type()
	if storedRaw, ok := rawCache.Load(objectType); ok {
		raw = storedRaw.(RawObject)
	} else {
		raw = NewSomeRawObjectWithTag(&objectValue, tagName)
		rawCache.Store(objectType, raw)
	}
	return raw
}
```