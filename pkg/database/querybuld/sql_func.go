package qb

import (
	"fmt"
	"strings"
)

func ASC(fieldName string) string {
	return fmt.Sprintf("%s ASC", fieldName)
}

func DESC(fieldName string) string {
	return fmt.Sprintf("%s DESC", fieldName)
}

func NullsFirst(fieldName string) string {
	return fmt.Sprintf("%s NULLS FIRST", fieldName)
}

func NullsLast(fieldName string) string {
	return fmt.Sprintf("%s NULLS LAST", fieldName)
}

func Rand() string {
	return "RAND()"
}

func Now() string {
	return "NOW()"
}

func CurrentDate() string {
	return "CURRENT_DATE()"
}

func CurrentTime() string {
	return "CURRENT_TIME()"
}

func Year(value string) string {
	return fmt.Sprintf("YEAR(%s)", value)
}

func Month(value string) string {
	return fmt.Sprintf("MONTH(%s)", value)
}

func Day(value string) string {
	return fmt.Sprintf("DAY(%s)", value)
}

func Lenght(value string) string {
	return fmt.Sprintf("LENGTH(%s)", value)
}

func CharLenght(value string) string {
	return fmt.Sprintf("CHAR_LENGTH(%s)", value)
}

func Upper(value string) string {
	return fmt.Sprintf("UPPER(%s)", value)
}

func Lower(value string) string {
	return fmt.Sprintf("LOWER(%s)", value)
}

func Concat(values ...string) string {
	return fmt.Sprintf("CONCAT(%s)", strings.Join(values, ", "))
}

func Substring(value string, start int, lenght int) string {
	if lenght > 0 {
		return fmt.Sprintf("SUBSTRING(%s, %v, %v)", value, start, lenght)
	} else {
		return fmt.Sprintf("CONCAT(%s, %v)", value, start)
	}
}

func Field(values ...any) string {
	res := "FIELD("
	for i := 0; i < len(values); i++ {
		if len(values)-1 == i {
			res += fmt.Sprintf("%v", values[i])
		} else {
			res += fmt.Sprintf("%v, ", values[i])
		}
	}
	res += ")"
	return res
}

func Distinct(values ...string) string {
	res := "DICTINCT "
	res += strings.Join(values, ", ")
	return res
}

func Using(name string, usingName string) string {
	return fmt.Sprintf("%s USING(%s)", name, usingName)
}

func Count(field string) string {
	return fmt.Sprintf("COUNT(%s)", field)
}

func AVG(field string) string {
	return fmt.Sprintf("AVG(%s)", field)
}

func Min(field string) string {
	return fmt.Sprintf("MIN(%s)", field)
}

func Max(field string) string {
	return fmt.Sprintf("MAX(%s)", field)
}

func GroupConcat(field string) string {
	return fmt.Sprintf("GROUP_CONCAT(%s)", field)
}
