package dbutils

import (
	"fmt"
)

type WHValue map[string]interface{}
type WHSliceValue map[string][]interface{}

type WHOutput struct {
	QueryStr  string
	QueryArgs []interface{}
}

func WHEquals(v WHValue, conjunction string) WHOutput {
	var queryArgs []interface{}
	var queryStr string
	i := -1
	for name, value := range v {
		i++
		queryArgs = append(queryArgs, value)
		if i == len(v)-1 {
			queryStr += fmt.Sprintf("%s = ?", name)
		} else {
			queryStr += fmt.Sprintf("%s = ? %s ", name, conjunction)
		}
	}
	return WHOutput{
		QueryStr:  queryStr,
		QueryArgs: queryArgs,
	}
}

func WHNotEquals(v WHValue, conjunction string) WHOutput {
	var queryArgs []interface{}
	var queryStr string
	i := -1
	for name, value := range v {
		i++
		queryArgs = append(queryArgs, value)
		if i == len(v)-1 {
			queryStr += fmt.Sprintf("%s != ?", name)
		} else {
			queryStr += fmt.Sprintf("%s != ? %s ", name, conjunction)
		}
	}
	return WHOutput{
		QueryStr:  queryStr,
		QueryArgs: queryArgs,
	}
}

func WHInSlice(v WHSliceValue, conjunction string) WHOutput {
	var queryStr string
	var queryArgs []interface{}
	i := -1
	for name, value := range v {
		i++
		var values string
		for i := 0; i < len(value); i++ {
			queryArgs = append(queryArgs, value[i])
			if i == len(value)-1 {
				values += "?"
			} else {
				values += "?, "
			}
		}
		if i == len(v)-1 {
			queryStr += fmt.Sprintf("%s IN %s", name, "("+values+")")
		} else {
			queryStr += fmt.Sprintf("%s IN %s %s ", name, "("+values+")", conjunction)
		}
	}
	return WHOutput{
		QueryStr:  queryStr,
		QueryArgs: queryArgs,
	}
}

func WHNotInSlice(v WHSliceValue, conjunction string) WHOutput {
	var queryStr string
	var queryArgs []interface{}
	i := -1
	for name, value := range v {
		i++
		var values string
		for i := 0; i < len(value); i++ {
			queryArgs = append(queryArgs, value[i])
			if i == len(value)-1 {
				values += "?"
			} else {
				values += "?, "
			}
		}
		if i == len(v)-1 {
			queryStr += fmt.Sprintf("%s NOT IN %s", name, "("+values+")")
		} else {
			queryStr += fmt.Sprintf("%s NOT IN %s %s ", name, "("+values+")", conjunction)
		}
	}
	return WHOutput{
		QueryStr:  queryStr,
		QueryArgs: queryArgs,
	}
}
