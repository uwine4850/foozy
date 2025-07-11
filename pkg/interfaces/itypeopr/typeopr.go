package itypeopr

// NewInstance creates a new instance of an object.
// The peculiarity of creating a new instance is that
// this method can be used from an already created object,
// which means that using this method you can transfer field
// values ​​and various settings to a new instance.
type NewInstance interface {
	New() (interface{}, error)
}
