package fstring

import (
	"unicode"
)

func ToLower(value string) string {
	for i, v := range value {
		return string(unicode.ToLower(v)) + value[i+1:]
	}
	return ""
}
