package mapper

import (
	"reflect"
)

type IMapper interface {
	Fill() error
}

type RawStruct interface {
	Type() reflect.Type
	Fields() *map[string]reflect.StructField
}
