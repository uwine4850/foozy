package utils

import (
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

func SliceContains[T comparable](slice []T, item T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return true
		}
	}
	return false
}
