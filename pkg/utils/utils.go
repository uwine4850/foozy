package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// SplitUrl separates the url by the "/" character. Skips empty slice values.
func SplitUrl(url string) []string {
	sp := strings.Split(url, "/")
	var res []string
	for i := 0; i < len(sp); i++ {
		if sp[i] == "" {
			continue
		}
		res = append(res, sp[i])
	}
	return res
}

// SplitUrlFromFirstSlug returns the left side of the url before the "<" sign.
func SplitUrlFromFirstSlug(url string) string {
	index := strings.Index(url, "<")
	if index == -1 {
		return url
	}
	return url[:index]
}

// SliceContains checks to see if the slice contains a value.
func SliceContains[T comparable](slice []T, item T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return true
		}
	}
	return false
}

// GenerateCsrfToken generates a CSRF token.
func GenerateCsrfToken() string {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic(err)
	}

	csrfToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return csrfToken
}

// MergeMap merges two maps into one.
// For example, if you pass Map1 and Map2, Map2 data will be added to Map1.
func MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2) {
	for key, value := range map2 {
		(*map1)[key] = value
	}
}

// Join outputs the slice in string format with the specified delimiter.
func Join[T any](elems []T, sep string) string {
	var res string
	for i := 0; i < len(elems); i++ {
		if i == len(elems)-1 {
			res += fmt.Sprintf("%v", elems[i])
		} else {
			res += fmt.Sprintf("%v%s ", elems[i], sep)
		}
	}
	return res
}
