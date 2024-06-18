package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/uwine4850/foozy/pkg/ferrors"
	"github.com/uwine4850/foozy/pkg/interfaces"
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
func GenerateCsrfToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	csrfToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return csrfToken, nil
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
	var res strings.Builder
	for i := 0; i < len(elems); i++ {
		if i == len(elems)-1 {
			res.WriteString(fmt.Sprintf("%v", elems[i]))
		} else {
			res.WriteString(fmt.Sprintf("%v%s ", elems[i], sep))
		}
	}
	return res.String()
}

func IsPointer(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Pointer
}

// CreateNewInstance сreates a new instance of a structure. The structure must implement the interface interfaces.INewInstance.
// The <new> argument takes a pointer to the structure that will contain the new instance.
func CreateNewInstance(ins interfaces.INewInstance, new interface{}) error {
	if !IsPointer(new) {
		panic(ferrors.ErrParameterNotPointer{Param: "new"})
	}
	var typeIns reflect.Type
	if reflect.TypeOf(ins).Kind() == reflect.Ptr {
		typeIns = reflect.TypeOf(ins).Elem()
	} else {
		typeIns = reflect.TypeOf(ins)
	}

	reflectIns := reflect.New(typeIns).Interface().(interfaces.INewInstance)
	newIns, err := reflectIns.New()
	if err != nil {
		return err
	}
	newInsInterface := reflect.ValueOf(newIns).Interface()
	reflect.ValueOf(new).Elem().Set(reflect.ValueOf(newInsInterface))
	return nil
}
