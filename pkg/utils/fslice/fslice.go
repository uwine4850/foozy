package fslice

// SliceContains checks to see if the slice contains a value.
func SliceContains[T comparable](slice []T, item T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return true
		}
	}
	return false
}

// AllStringItemsEmpty checks if all items in the slice are empty strings.
func AllStringItemsEmpty(slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] != "" {
			return false
		}
	}
	return true
}
