package itypeopr

type INewInstance interface {
	New() (interface{}, error)
}

type IPtr interface {
	New(value interface{}) IPtr
	Ptr() interface{}
}
