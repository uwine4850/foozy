package qb

import "fmt"

// Comparison operator. It is used to compare values in the Compare structure.
type CompareOperator string

const (
	EQUAL      CompareOperator = "="
	NOT_EQUAL  CompareOperator = "!="
	GREATER    CompareOperator = ">"
	LESS       CompareOperator = "<"
	IS         CompareOperator = "IS"
	IS_NOT     CompareOperator = "IS NOT"
	IN         CompareOperator = "IN"
	NOT_IN     CompareOperator = "NOT IN"
	LIKE       CompareOperator = "LIKE"
	NOT_LIKE   CompareOperator = "NOT LIKE"
	REGEXP     CompareOperator = "REGEXP"
	NOT_REGEXP CompareOperator = "NOT REGEXP"
)

// Combines two conditions.
type Uniteers string

// IsUniteers checks whether the value is of type Uniteers.
func IsUniteers(value any) bool {
	if _, ok := value.(Uniteers); ok {
		return true
	}
	return false
}

const (
	AND Uniteers = "AND"
	OR  Uniteers = "OR"
)

// Special data types that do not belong to a standard type.
type SpecialType string

// IsSpecialType checks whether the value is of type SpecialType.
func IsSpecialType(value any) bool {
	if _, ok := value.(SpecialType); ok {
		return true
	}
	return false
}

const (
	NULL  SpecialType = "NULL"
	TRUE  SpecialType = "True"
	FALSE SpecialType = "False"
)

// Sql standard field data type.
type T struct {
	value string
}

// Value string value of the selected data type.
func (t T) Value() string {
	return t.value
}

func (t T) Text() T {
	t.value = "TEXT"
	return t
}

func (t T) Int(lenght int) T {
	t.value = fmt.Sprintf("INT(%v)", lenght)
	return t
}

func (t T) Varchar(lenght int) T {
	t.value = fmt.Sprintf("VARCHAR(%v)", lenght)
	return t
}

func (t T) Date() T {
	t.value = "DATE"
	return t
}

func (t T) Boolean() T {
	t.value = "BOOLEAN"
	return t
}
