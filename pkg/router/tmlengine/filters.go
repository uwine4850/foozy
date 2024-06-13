package tmlengine

import (
	"fmt"
	"github.com/flosch/pongo2"
	"html"
	"strings"
)

type Filter struct {
	Name string
	Fn   pongo2.FilterFunction
}

var BuiltinFilters = []Filter{
	{
		// Converts a string with escaped characters to just a string.
		"unescape",
		func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
			return pongo2.AsValue(html.UnescapeString(in.String())), nil
		},
	},
	{
		// Converts a slice of any type to a string.
		"strslice",
		func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
			var strSlice string
			if in.CanSlice() {
				inSlice := in
				var newSlice []string
				for i := 0; i < inSlice.Len(); i++ {
					newSlice = append(newSlice, fmt.Sprintf("%v", inSlice.Index(i).Interface()))
				}
				strSlice = strings.Join(newSlice, ", ")
			}
			return pongo2.AsValue(strSlice), nil
		},
	},
}

// RegisterGlobalFilter globally registers the pongo2 filter.
func RegisterGlobalFilter(name string, fn pongo2.FilterFunction) error {
	err := pongo2.RegisterFilter(name, fn)
	if err != nil {
		return err
	}
	return nil
}

// RegisterMultipleGlobalFilter globally registers multiple pongo2 filters.
func RegisterMultipleGlobalFilter(filters []Filter) error {
	for i := 0; i < len(filters); i++ {
		err := RegisterGlobalFilter(filters[i].Name, filters[i].Fn)
		if err != nil {
			return err
		}
	}
	return nil
}
