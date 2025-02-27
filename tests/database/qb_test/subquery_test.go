package qbtest

import (
	"reflect"
	"testing"

	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
)

func TestSubquery(t *testing.T) {
	sq := qb.SQ(true, qb.NewNoDbQB().SelectFrom(qb.Count("*"), "tableName"))
	q := syncQB().Select(sq).As("count")
	q.Merge()
	if q.String() != "SELECT (SELECT COUNT(*) FROM tableName) AS count" {
		t.Error(testErrorText)
	}
}

func TestParseSubQuery(t *testing.T) {
	var qString string
	var qArgs []any
	sq := qb.SQ(false, qb.NewNoDbQB().SelectFrom(qb.Count("*"), "tableName").Where(qb.Compare("id", qb.NOT_EQUAL, 1)))
	qb.ParseSubQuery(sq, &qString, &qArgs)
	if qString != "SELECT COUNT(*) FROM tableName WHERE id != ?" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(qArgs, []any{1}) {
		t.Error(testErrorArgs)
	}
}

func TestIsSubQuery(t *testing.T) {
	sq := qb.SQ(false, qb.NewNoDbQB())
	if !qb.IsSubQuery(sq) {
		t.Error("subquery definition error")
	}
}
