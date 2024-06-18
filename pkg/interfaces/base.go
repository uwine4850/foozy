package interfaces

type INewInstance interface {
	New() (interface{}, error)
}
