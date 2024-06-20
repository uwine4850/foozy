package fmap

// MergeMap merges two maps into one.
// For example, if you pass Map1 and Map2, Map2 data will be added to Map1.
func MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2) {
	for key, value := range map2 {
		(*map1)[key] = value
	}
}
