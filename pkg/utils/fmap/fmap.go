package fmap

import (
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

// MergeMap merges two maps into one.
// For example, if you pass Map1 and Map2, Map2 data will be added to Map1.
func MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2) {
	for key, value := range map2 {
		(*map1)[key] = value
	}
}

// Compare map values. It is important that the keys and values ​​match.
// exclude - keys that do not need to be taken into account.
func Compare[T1 comparable, T2 comparable](map1 *map[T1]T2, map2 *map[T1]T2, exclude []T1) bool {
	for key, value := range *map1 {
		if fslice.SliceContains(exclude, key) {
			continue
		}
		value2, ok := (*map2)[key]
		if !ok {
			return false
		} else {
			if value != value2 {
				return false
			}
		}
	}
	return true
}
