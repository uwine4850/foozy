package fstring

import (
	"fmt"
	"strings"
)

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

// Join outputs the slice in string format with the specified delimiter.
func Join[T any](elems []T, sep string) string {
	var res strings.Builder
	for i := 0; i < len(elems); i++ {
		if i == len(elems)-1 {
			res.WriteString(fmt.Sprintf("%v", elems[i]))
		} else {
			res.WriteString(fmt.Sprintf("%v%s", elems[i], sep))
		}
	}
	return res.String()
}
