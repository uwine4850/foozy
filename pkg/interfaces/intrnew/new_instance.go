package intrnew

type INewInstance interface {
	New() (interface{}, error)
}
