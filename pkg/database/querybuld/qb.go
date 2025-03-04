package qb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

// Col structure that implements a sql column in a table.
type Col struct {
	Name    string
	Type    T
	Null    bool
	AI      bool
	Default string
	PK      bool
}

// String the string value of a sql column in a table.
func (c *Col) String() string {
	var null string
	if c.Null {
		null = "NULL"
	} else {
		null = "NOT NULL"
	}
	res := fmt.Sprintf("%s %s %s", c.Name, c.Type.Value(), null)
	if c.AI {
		res += " AUTO_INCREMENT"
	}
	if c.Default != "" {
		res += " " + c.Default
	}
	if c.PK {
		res += " PRIMARY KEY"
	}
	return res
}

// TableQB query build sql table.
// Creates a table in the database if it does not exist.
type TableQB struct {
	parts       []string
	queryString string
}

func NewTableQB() *TableQB {
	return &TableQB{}
}

// Create creates a table.
func (tqb *TableQB) Create(tableName string, cols []Col) *TableQB {
	var colsString string
	for i := 0; i < len(cols); i++ {
		if len(cols)-1 == i {
			colsString += cols[i].String() + " "
		} else {
			colsString += cols[i].String() + ", "
		}
	}
	qString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, colsString)
	tqb.parts = append(tqb.parts, qString)
	return tqb
}

// Drop executes the DROP command for tables.
func (tqb *TableQB) Drop(tables ...string) *TableQB {
	tqb.parts = append(tqb.parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", strings.Join(tables, ", ")))
	return tqb
}

// Drop executes the TRUNCATE command for tables.
func (tqb *TableQB) Truncate(tables ...string) *TableQB {
	tqb.parts = append(tqb.parts, fmt.Sprintf("TRUNCATE TABLE IF EXISTS %s;", strings.Join(tables, ", ")))
	return tqb
}

// FK creates a foreign key relationship between tables.
func (tqb *TableQB) FK(tableName string, constraintName string, fkCol string, refCol string, onUpdate string, onDelete string) *TableQB {
	var _onUpdate string
	if onUpdate != "" {
		_onUpdate = "ON UPDATE " + onUpdate
	}
	var _onDelete string
	if onDelete != "" {
		_onDelete = "ON DELETE " + onDelete
	}
	q := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s %s %s;", tableName, constraintName, fkCol, refCol, _onDelete, _onUpdate)
	tqb.parts = append(tqb.parts, q)
	return tqb
}

// String value of sql command to create table.
func (tqb *TableQB) String() string {
	return tqb.queryString
}

func (tqb *TableQB) Build() *TableQB {
	tqb.queryString = strings.Join(tqb.parts, " ")
	return tqb
}

// QB query build structure for creating sql queries.
// This structure consists of more narrowly focused objects that implement certain sql functionality.
type QB struct {
	dataOperationQB
	filterQB
	joinQB
	customFunc
	queryParts  []string
	queryString string
	queryArgs   []any
	syncQ       interfaces.ISyncQueries
	asyncQ      interfaces.IAsyncQueries
	asyncKey    string
}

func NewNoDbQB() *QB {
	qb := &QB{}
	qb.dataOperationQB.SetQB(qb)
	qb.filterQB.SetQB(qb)
	qb.joinQB.SetQB(qb)
	qb.customFunc.SetQB(qb)
	return qb
}

func NewSyncQB(syncQ interfaces.ISyncQueries) *QB {
	qb := &QB{
		syncQ: syncQ,
	}
	qb.dataOperationQB.SetQB(qb)
	qb.filterQB.SetQB(qb)
	qb.joinQB.SetQB(qb)
	qb.customFunc.SetQB(qb)
	return qb
}

func NewAsyncQB(asyncQ interfaces.IAsyncQueries, key string) *QB {
	qb := &QB{
		asyncQ:   asyncQ,
		asyncKey: key,
	}
	qb.dataOperationQB.SetQB(qb)
	qb.filterQB.SetQB(qb)
	qb.joinQB.SetQB(qb)
	qb.customFunc.SetQB(qb)
	return qb
}

// AppendPart adds a part of the sql query to the overall slice.
func (qb *QB) AppendPart(part string) {
	qb.queryParts = append(qb.queryParts, part)
}

// AppendArgs adds positional sql arguments to a common argument slice.
func (qb *QB) AppendArgs(args []any) {
	qb.queryArgs = append(qb.queryArgs, args...)
}

// Merge uses the string parts of the sql query that were added to [queryParts]
// to combine them into a coherent sql query.
func (qb *QB) Merge() {
	qb.queryString = strings.Join(qb.queryParts, " ")
}

// String outputs sql string created with Build function.
func (qb *QB) String() string {
	return qb.queryString
}

// Args outputs positional sql arguments.
func (qb *QB) Args() []any {
	return qb.queryArgs
}

func (qb *QB) Query() ([]map[string]interface{}, error) {
	qb.Merge()
	if qb.syncQ != nil {
		return qb.syncQ.Query(qb.String(), qb.Args()...)
	}
	if qb.asyncQ != nil {
		qb.asyncQ.Query(qb.asyncKey, qb.String(), qb.Args()...)
		return nil, nil
	}
	return nil, errors.New("no synchronous or asynchronous QB handler found")
}

func (qb *QB) Exec() (map[string]interface{}, error) {
	qb.Merge()
	if qb.syncQ != nil {
		return qb.syncQ.Exec(qb.String(), qb.Args()...)
	}
	if qb.asyncQ != nil {
		qb.asyncQ.Exec(qb.asyncKey, qb.String(), qb.Args()...)
		return nil, nil
	}
	return nil, errors.New("no synchronous or asynchronous QB handler found")
}

// mergeTwoQB merges two QB instances into one.
func mergeTwoQB(qb1 *QB, qb2 *QB) *QB {
	var newQB *QB
	if qb1.syncQ != nil {
		newQB = NewSyncQB(qb1.syncQ)
	}
	if qb1.asyncQ != nil {
		newQB = NewAsyncQB(qb1.asyncQ, qb1.asyncKey)
	}
	qb1.Merge()
	qb2.Merge()
	return newQB
}

// Union adds a UNION command for two sql queries passed using QB.
func Union(qb1 *QB, qb2 *QB) *QB {
	newQB := mergeTwoQB(qb1, qb2)
	newQB.AppendPart(fmt.Sprintf("(%s)", qb1.String()))
	newQB.AppendArgs(qb1.Args())
	newQB.AppendPart("UNION")
	newQB.AppendPart(fmt.Sprintf("(%s)", qb2.String()))
	newQB.AppendArgs(qb2.Args())
	return newQB
}

// UnionAll adds a UNION ALL command for two sql queries passed using QB.
func UnionAll(qb1 *QB, qb2 *QB) *QB {
	newQB := mergeTwoQB(qb1, qb2)
	newQB.AppendPart(fmt.Sprintf("(%s)", qb1.String()))
	newQB.AppendArgs(qb1.Args())
	newQB.AppendPart("UNION ALL")
	newQB.AppendPart(fmt.Sprintf("(%s)", qb2.String()))
	newQB.AppendArgs(qb2.Args())
	return newQB
}

// Intersect adds a INTERSECT command for two sql queries passed using QB.
func Intersect(qb1 *QB, qb2 *QB) *QB {
	newQB := mergeTwoQB(qb1, qb2)
	newQB.AppendPart(fmt.Sprintf("(%s)", qb1.String()))
	newQB.AppendArgs(qb1.Args())
	newQB.AppendPart("INTERSECT")
	newQB.AppendPart(fmt.Sprintf("(%s)", qb2.String()))
	newQB.AppendArgs(qb2.Args())
	return newQB
}

// Except adds a EXCEPT command for two sql queries passed using QB.
func Except(qb1 *QB, qb2 *QB) *QB {
	newQB := mergeTwoQB(qb1, qb2)
	newQB.AppendPart(fmt.Sprintf("(%s)", qb1.String()))
	newQB.AppendArgs(qb1.Args())
	newQB.AppendPart("EXCEPT")
	newQB.AppendPart(fmt.Sprintf("(%s)", qb2.String()))
	newQB.AppendArgs(qb2.Args())
	return newQB
}

// dataOperationQB object for operations with table data.
// Metgod name corresponds to the corresponding sql operation.
type dataOperationQB struct {
	qb *QB
}

func (doQB *dataOperationQB) SetQB(qb *QB) {
	doQB.qb = qb
}

func (doQB *dataOperationQB) Custom(value string, args ...any) *QB {
	doQB.qb.AppendPart(value)
	doQB.qb.AppendArgs(args)
	return doQB.qb
}

func (doQB *dataOperationQB) Select(values ...any) *QB {
	var qString string
	qArgs := []any{}
	for i := 0; i < len(values); i++ {
		if !ParseSubQuery(values[i], &qString, &qArgs) && !processingConditionBuilder(values[i], &qString, &qArgs) {
			if reflect.TypeOf(values[i]).Kind() == reflect.String {
				qString += " " + strings.Trim(values[i].(string), " ")
			} else {
				panic(fmt.Sprintf("%s data type is not supported", reflect.TypeOf(values[i])))
			}
		}
	}
	doQB.qb.AppendPart(fmt.Sprintf("SELECT %s", strings.Trim(qString, " ")))
	doQB.qb.AppendArgs(qArgs)
	return doQB.qb
}

func (doQB *dataOperationQB) As(name string) *QB {
	doQB.qb.AppendPart(fmt.Sprintf("AS %s", name))
	return doQB.qb
}

func (doQB *dataOperationQB) SelectFrom(target any, from any) *QB {
	var targetValue any
	var fromValue any
	qArgs := []any{}
	var sqQString string
	if ok := ParseSubQuery(target, &sqQString, &qArgs); ok {
		targetValue = sqQString
	} else {
		targetValue = target
	}
	if ok := ParseSubQuery(from, &sqQString, &qArgs); ok {
		fromValue = sqQString
	} else {
		fromValue = from
	}
	doQB.qb.AppendPart(fmt.Sprintf("SELECT %v FROM %v", targetValue, fromValue))
	doQB.qb.AppendArgs(qArgs)
	return doQB.qb
}

func (doQB *dataOperationQB) Insert(tableName string, params map[string]any) *QB {
	keys, values := dbutils.ParseParams(params)
	doQB.qb.AppendPart(fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ( %s )", tableName, strings.Join(keys, ", "), dbutils.RepeatValues(len(values), ",")))
	doQB.qb.AppendArgs(values)
	return doQB.qb
}

func (doQB *dataOperationQB) Update(tableName string, params map[string]any) *QB {
	strVal, args := dbutils.ParseMapAsEquals(&params)
	doQB.qb.AppendPart(fmt.Sprintf("UPDATE %s SET %s", tableName, strVal))
	doQB.qb.AppendArgs(args)
	return doQB.qb
}

func (doQB *dataOperationQB) Delete(tableName string) *QB {
	doQB.qb.AppendPart(fmt.Sprintf("DELETE FROM %s", tableName))
	return doQB.qb
}

// joinQB operations with joining table values.
// Metgod name corresponds to the corresponding sql operation.
type joinQB struct {
	qb *QB
}

func (jqb *joinQB) SetQB(qb *QB) {
	jqb.qb = qb
}

func (jqb *joinQB) InnerJoin(tableName string, values ...any) *QB {
	var qString string
	qArgs := []any{}
	if len(values) != 0 {
		processingConditionValues(values, &qString, &qArgs)
		jqb.qb.AppendPart(fmt.Sprintf("INNER JOIN %s ON %s", tableName, qString))
	} else {
		jqb.qb.AppendPart(fmt.Sprintf("INNER JOIN %s", tableName))
	}
	jqb.qb.AppendArgs(qArgs)
	return jqb.qb
}

func (jqb *joinQB) LeftJoin(tableName string, values ...any) *QB {
	var qString string
	qArgs := []any{}
	if len(values) != 0 {
		processingConditionValues(values, &qString, &qArgs)
		jqb.qb.AppendPart(fmt.Sprintf("LEFT JOIN %s ON %s", tableName, qString))
	} else {
		jqb.qb.AppendPart(fmt.Sprintf("LEFT JOIN %s", tableName))
	}
	jqb.qb.AppendArgs(qArgs)
	return jqb.qb
}

func (jqb *joinQB) RightJoin(tableName string, values ...any) *QB {
	var qString string
	qArgs := []any{}
	if len(values) != 0 {
		processingConditionValues(values, &qString, &qArgs)
		jqb.qb.AppendPart(fmt.Sprintf("RIGHT JOIN %s ON %s", tableName, qString))
	} else {
		jqb.qb.AppendPart(fmt.Sprintf("RIGHT JOIN %s", tableName))
	}
	jqb.qb.AppendArgs(qArgs)
	return jqb.qb
}

// The customFunc object adds the ability to use [subquery] anywhere in the query.
// This allows you to make more flexible queries.
type customFunc struct {
	qb *QB
}

func (cf *customFunc) SetQB(qb *QB) {
	cf.qb = qb
}

// Func runs a [subquery] in the selected [QB] query fragment.
func (cf *customFunc) Func(_subquery *subquery) *QB {
	var qString string
	qArgs := []any{}
	if ParseSubQuery(_subquery, &qString, &qArgs) {
		cf.qb.AppendPart(qString)
		cf.qb.AppendArgs(qArgs)
	}
	return cf.qb
}

// sql subquery that is inserted in the middle of the main query.
// The bracket field determines whether this query will be bracketed.
//
// It is important to specify that a new instance of the [QB] object is used.
// The [QB] object can also be without an initialized database, since it is not used.
// The [QB] structure is used only to receive the sql string and arguments, the query itself is not sent.
type subquery struct {
	qb      *QB
	bracket bool
}

func SQ(bracket bool, qb *QB) *subquery {
	qb.Merge()
	return &subquery{
		qb:      qb,
		bracket: bracket,
	}
}

func (sq *subquery) String() string {
	if sq.bracket {
		return fmt.Sprintf("(%s)", sq.qb.String())
	} else {
		return fmt.Sprintf("%s", sq.qb.String())
	}
}

func (sq *subquery) Args() []any {
	return sq.qb.Args()
}

// IsSubQuery checks if the value is of type subQuery.
func IsSubQuery(val any) bool {
	if reflect.TypeOf(val) == reflect.TypeOf(&subquery{}) {
		return true
	}
	return false
}

// ParseSubQuery processes an instance of subQuery.
// outQueryString — sql string of the subquery.
// outQueryArgs — position arguments of sql subquery.
func ParseSubQuery(sq any, outQueryString *string, outQueryArgs *[]any) bool {
	if IsSubQuery(sq) {
		_subQuery := sq.(*subquery)
		*outQueryString = _subQuery.String()
		*outQueryArgs = append(*outQueryArgs, _subQuery.Args()...)
		return true
	}
	return false
}
