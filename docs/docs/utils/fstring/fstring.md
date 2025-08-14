## fstring

#### ToLower
Converts the first letter of a string to lowercase.
```golang
func ToLower(value string) string {
	for i, v := range value {
		return string(unicode.ToLower(v)) + value[i+1:]
	}
	return ""
}
```