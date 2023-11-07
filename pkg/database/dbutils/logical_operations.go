package dbutils

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/utils"
)

type WHOutput struct {
	QueryStr  string
	QueryArgs []interface{}
}

func WHEquals(v map[string]interface{}, conjunction string) WHOutput {
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

func WHNotEquals(v map[string]interface{}, conjunction string) WHOutput {
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

func WHInSlice(v map[string][]interface{}, conjunction string) WHOutput {
	var queryStr string
	i := -1
	for name, value := range v {
		i++
		if i == len(v)-1 {
			queryStr += fmt.Sprintf("%s IN %s", name, "("+utils.Join(value, ",")+")")
		} else {
			queryStr += fmt.Sprintf("%s IN %s %s ", name, "("+utils.Join(value, ",")+")", conjunction)
		}
	}
	return WHOutput{
		QueryStr:  queryStr,
		QueryArgs: make([]interface{}, 0),
	}
}

func WHNotInSlice(v map[string][]interface{}, conjunction string) WHOutput {
	var queryStr string
	i := -1
	for name, value := range v {
		i++
		if i == len(v)-1 {
			queryStr += fmt.Sprintf("%s NOT IN %s", name, "("+utils.Join(value, ",")+")")
		} else {
			queryStr += fmt.Sprintf("%s NOT IN %s %s ", name, "("+utils.Join(value, ",")+")", conjunction)
		}
	}
	return WHOutput{
		QueryStr:  queryStr,
		QueryArgs: make([]interface{}, 0),
	}
}
