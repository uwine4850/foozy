## fslice

#### SliceContains
Checks to see if the slice contains a value.
```golang
func SliceContains[T comparable](slice []T, item T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return true
		}
	}
	return false
}
```

#### AllStringItemsEmpty
Checks if all items in the slice are empty strings.
```golang
func AllStringItemsEmpty(slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] != "" {
			return false
		}
	}
	return true
}
```