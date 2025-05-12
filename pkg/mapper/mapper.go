package mapper

import (
	"reflect"
)

type IMapper interface {
	Fill() error
}

// RawStruct is used to store data of type undefined structure.
//
// This structure is intended to optimize the work of the reflect package with the structure.
// The type of the structure itself and its field types are stored here. This object reduces
// the load on the structure parser, now it is not necessary to loop through all fields.
// Also this object is convenient because it can store only selected fields of the structure,
// which makes it more convenient to use.
//
// Fields are passed by the map so that it is possible to give them a custom tag name. If the tag
// is not needed, the key can be the field name.
//
// IMPORTANT: to get better performance you should store an instance of this object in a separate immutable variable.
type RawStruct interface {
	// Type returns the type of the stored object.
	Type() reflect.Type
	// Fields returns the selected fields of the stored object.
	// It can store not all fields, but only those added by the implementation.
	Fields() *map[string]reflect.StructField
}
