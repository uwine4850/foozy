package fstring

import (
	"strings"
	"unicode"
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

func ToLower(value string) string {
	for i, v := range value {
		return string(unicode.ToLower(v)) + value[i+1:]
	}
	return ""
}
