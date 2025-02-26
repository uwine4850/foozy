package qb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

type IConditionBuilder interface {
	Build()
	QString() string
	QArgs() []any
}

// processingConditionBuilder processes an instance of IConditionBuilder.
// outQueryString — sql string of the subquery.
// outQueryArgs — position arguments of sql subquery.
func processingConditionBuilder(builder any, outString *string, outArgs *[]any) bool {
	if !typeopr.IsPointer(builder) {
		return false
	}
	if typeopr.IsImplementInterface(typeopr.Ptr{}.New(builder), (*IConditionBuilder)(nil)) {
		conditionBuilder := builder.(IConditionBuilder)
		conditionBuilder.Build()
		*outString += conditionBuilder.QString()
		*outArgs = append(*outArgs, conditionBuilder.QArgs()...)
		return true
	}
	return false
}

// processingConditionValues processes the stratification values.
// The value can be a simple string or an IConditionBuilder instance.
// outQueryString — sql string of the subquery.
// outQueryArgs — position arguments of sql subquery.
func processingConditionValues(values []any, qString *string, qArgs *[]any) {
	for i := 0; i < len(values); i++ {
		if IsUniteers(values[i]) {
			*qString += " " + string(values[i].(Uniteers)) + " "
			continue
		}
		if !processingConditionBuilder(values[i], qString, qArgs) {
			if reflect.TypeOf(values[i]) == reflect.TypeOf("") {
				*qString += " " + values[i].(string)
			} else {
				panic(fmt.Sprintf("%s data type is not supported", reflect.TypeOf(values[i])))
			}
		}
	}
}

// filterQB query build for filtering.
// Metgod name corresponds to the corresponding sql operation.
type filterQB struct {
	mainQB *QB
}

func (fqb *filterQB) SetQB(qb *QB) {
	fqb.mainQB = qb
}

func (fqb *filterQB) Where(values ...any) *QB {
	qString := "WHERE "
	qArgs := []any{}
	processingConditionValues(values, &qString, &qArgs)
	fqb.mainQB.AppendPart(qString)
	fqb.mainQB.AppendArgs(qArgs)
	return fqb.mainQB
}

func (fqb *filterQB) OrderBy(conditions ...string) *QB {
	var qString string
	for i := 0; i < len(conditions); i++ {
		qString = strings.Join(conditions, ", ")
	}
	fqb.mainQB.AppendPart(fmt.Sprintf("ORDER BY %s", qString))
	return fqb.mainQB
}

func (fqb *filterQB) GroupBy(values ...string) *QB {
	qString := "GROUP BY "
	qString += strings.Join(values, ", ")
	fqb.mainQB.AppendPart(qString)
	return fqb.mainQB
}

func (fqb *filterQB) Having(values ...any) *QB {
	qString := "HAVING "
	qArgs := []any{}
	processingConditionValues(values, &qString, &qArgs)
	fqb.mainQB.AppendPart(qString)
	fqb.mainQB.AppendArgs(qArgs)
	return fqb.mainQB
}

func (fqb *filterQB) Limit(number int) *QB {
	fqb.mainQB.AppendPart(fmt.Sprintf("LIMIT %v", number))
	return fqb.mainQB
}

func (fqb *filterQB) Offset(number int) *QB {
	fqb.mainQB.AppendPart(fmt.Sprintf("OFFSET %v", number))
	return fqb.mainQB
}

// compare object that performs comparisons of two values.
type compare struct {
	leftOperand  any
	operator     CompareOperator
	rightOperand any
	leftValue    any
	rightValue   any
	args         []any
}

func Compare(leftOperand any, operator CompareOperator, rightOperand any) *compare {
	return &compare{
		leftOperand:  leftOperand,
		operator:     operator,
		rightOperand: rightOperand,
	}
}

func (c *compare) Build() {
	var leftOperandSQValue string
	var rightOperandSQValue string
	if !processingConditionBuilder(c.leftOperand, &leftOperandSQValue, &c.args) {
		ParseSubQuery(c.leftOperand, &leftOperandSQValue, &c.args)
	}
	if IsSpecialType(c.rightOperand) {
		rightOperandSQValue = string(c.rightOperand.(SpecialType))
	} else {
		if !processingConditionBuilder(c.rightOperand, &rightOperandSQValue, &c.args) {
			if !ParseSubQuery(c.rightOperand, &rightOperandSQValue, &c.args) {
				c.args = append(c.args, c.rightOperand)
			}
		}
	}

	if leftOperandSQValue != "" {
		c.leftValue = leftOperandSQValue
	} else {
		c.leftValue = c.leftOperand
	}
	if rightOperandSQValue != "" {
		c.rightValue = rightOperandSQValue
	} else {
		c.rightValue = "?"
	}
}

func (c *compare) QString() string {
	return fmt.Sprintf("%v %v %v", c.leftValue, c.operator, c.rightValue)
}

func (c *compare) QArgs() []any {
	return c.args
}

// noArgsCompare does almost everything that a normal [compare] object does.
// The peculiarity is that this object does not pass the right operand as a
// positional argument. So the "?" sign will not be used.
// It is important to clarify that the structure still processes and
// passes arguments of external objects like [subQuery].
type noArgsCompare struct {
	leftOperand  any
	operator     CompareOperator
	rightOperand any
	leftValue    any
	rightValue   any
	args         []any
}

func NoArgsCompare(leftOperand any, operator CompareOperator, rightOperand any) *noArgsCompare {
	return &noArgsCompare{
		leftOperand:  leftOperand,
		operator:     operator,
		rightOperand: rightOperand,
	}
}

func (c *noArgsCompare) Build() {
	var leftOperandSQValue string
	var rightOperandSQValue string
	if !processingConditionBuilder(c.leftOperand, &leftOperandSQValue, &c.args) {
		ParseSubQuery(c.leftOperand, &leftOperandSQValue, &c.args)
	}
	if IsSpecialType(c.rightOperand) {
		rightOperandSQValue = string(c.rightOperand.(SpecialType))
	} else {
		if !processingConditionBuilder(c.rightOperand, &rightOperandSQValue, &c.args) {
			ParseSubQuery(c.rightOperand, &rightOperandSQValue, &c.args)
		}
	}

	if leftOperandSQValue != "" {
		c.leftValue = leftOperandSQValue
	} else {
		c.leftValue = c.leftOperand
	}
	if rightOperandSQValue != "" {
		c.rightValue = rightOperandSQValue
	} else {
		c.rightValue = c.rightOperand
	}
}

func (c *noArgsCompare) QString() string {
	return fmt.Sprintf("%v %v %v", c.leftValue, c.operator, c.rightValue)
}

func (c *noArgsCompare) QArgs() []any {
	return c.args
}

// between processes the sql command BETWEEN.
type between struct {
	columnName   string
	leftOperand  any
	rightOperand any
	leftValue    any
	rightValue   any
	args         []any
}

func Between(columnName string, leftOperand any, rightOperand any) *between {
	return &between{
		columnName:   columnName,
		leftOperand:  leftOperand,
		rightOperand: rightOperand,
	}
}

func (b *between) Build() {
	var leftOperandSQValue string
	var rightOperandSQValue string
	if !ParseSubQuery(b.leftOperand, &leftOperandSQValue, &b.args) {
		b.args = append(b.args, b.leftOperand)
	}
	if !ParseSubQuery(b.rightOperand, &rightOperandSQValue, &b.args) {
		b.args = append(b.args, b.rightOperand)
	}

	if leftOperandSQValue != "" {
		b.leftValue = leftOperandSQValue
	} else {
		b.leftValue = "?"
	}
	if rightOperandSQValue != "" {
		b.rightValue = rightOperandSQValue
	} else {
		b.rightValue = "?"
	}
}

func (b *between) QString() string {
	return fmt.Sprintf("%s BETWEEN %v AND %v", b.columnName, b.leftValue, b.rightValue)
}

func (b *between) QArgs() []any {
	return b.args
}

// notBetween processes the sql command NOT BETWEEN.
type notBetween struct {
	between
}

func NotBetween(columnName string, leftOperand any, rightOperand any) *notBetween {
	return &notBetween{
		between{
			columnName:   columnName,
			leftOperand:  leftOperand,
			rightOperand: rightOperand,
		},
	}
}

func (nb *notBetween) QString() string {
	return fmt.Sprintf("%s NOT BETWEEN %v AND %v", nb.columnName, nb.leftValue, nb.rightValue)
}

// array processes the data as a sql list.
type array struct {
	values  []any
	qString string
	qArgs   []any
}

func Array(values ...any) *array {
	return &array{
		values: values,
	}
}

func (arr *array) Build() {
	arr.qString = "("
	for i := 0; i < len(arr.values); i++ {
		var q string
		if processingConditionBuilder(arr.values[i], &q, &arr.qArgs) {
			arr.qString += " " + q
			continue
		}
		arr.qArgs = append(arr.qArgs, arr.values[i])
		arr.qString += " ?"
	}
	arr.qString += " )"
}

func (arr *array) QString() string {
	return arr.qString
}

func (arr *array) QArgs() []any {
	return arr.qArgs
}

// exists processes the sql command EXISTS.
type exists struct {
	sq      *subquery
	qString string
	qArgs   []any
}

func Exists(sq *subquery) *exists {
	return &exists{
		sq: sq,
	}
}

func (e *exists) Build() {
	ParseSubQuery(e.sq, &e.qString, &e.qArgs)
}

func (e *exists) QString() string {
	return fmt.Sprintf("EXISTS(%s)", e.qString)
}

func (e *exists) QArgs() []any {
	return e.qArgs
}
