package qbtest

import (
	"reflect"
	"testing"

	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
)

func TestCompare(t *testing.T) {
	comp := qb.Compare("id", qb.EQUAL, 1)
	comp.Build()
	if comp.QString() != "id = ?" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(comp.QArgs(), []any{1}) {
		t.Error(testErrorArgs)
	}
}

func TestCompareSQ(t *testing.T) {
	sq := qb.SQ(true, qb.NewNoDbQB().SelectFrom(qb.Count("*"), "tableName").Where(qb.Compare("id", qb.NOT_EQUAL, 1)))
	comp := qb.Compare("id", qb.EQUAL, sq)
	comp.Build()
	if comp.QString() != "id = (SELECT COUNT(*) FROM tableName WHERE id != ?)" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(comp.QArgs(), []any{1}) {
		t.Error(testErrorArgs)
	}
}

func TestBetween(t *testing.T) {
	bet := qb.Between("id", 1, 5)
	bet.Build()
	if bet.QString() != "id BETWEEN ? AND ?" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(bet.QArgs(), []any{1, 5}) {
		t.Error(testErrorArgs)
	}
}

func TestBetweenSQ(t *testing.T) {
	sq := qb.SQ(true, qb.NewNoDbQB().SelectFrom(qb.Count("*"), "tableName").Where(qb.Compare("id", qb.NOT_EQUAL, 1)))
	bet := qb.Between("id", sq, sq)
	bet.Build()
	if bet.QString() != "id BETWEEN (SELECT COUNT(*) FROM tableName WHERE id != ?) AND (SELECT COUNT(*) FROM tableName WHERE id != ?)" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(bet.QArgs(), []any{1, 1}) {
		t.Error(testErrorArgs)
	}
}

func TestExists(t *testing.T) {
	sq := qb.SQ(false, qb.NewNoDbQB().SelectFrom(1, "tableName").Where(qb.Compare("id", qb.EQUAL, 1)))
	ex := qb.Exists(sq)
	ex.Build()
	if ex.QString() != "EXISTS(SELECT 1 FROM tableName WHERE id = ?)" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(ex.QArgs(), []any{1}) {
		t.Error(testErrorArgs)
	}
}

func TestWhere(t *testing.T) {
	q := syncQB.SelectFrom("id", "table").Where(qb.Compare("id", qb.EQUAL, 1))
	q.Merge()
	if q.String() != "SELECT id FROM table WHERE id = ?" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(q.Args(), []any{1}) {
		t.Error(testErrorArgs)
	}
}

func TestWhereInArray(t *testing.T) {
	q := syncQB.SelectFrom("id", "table").Where(qb.Compare("id", qb.IN, qb.Array(1, 3, 4)))
	q.Merge()
	if q.String() != "SELECT id FROM table WHERE id IN ( ? ? ? )" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(q.Args(), []any{1, 3, 4}) {
		t.Error(testErrorArgs)
	}
}

func TestOrderBy(t *testing.T) {
	q := syncQB.SelectFrom("id", "table").OrderBy(qb.DESC("id"))
	q.Merge()
	if q.String() != "SELECT id FROM table ORDER BY id DESC" {
		t.Error(testErrorText)
	}
}

func TestGroupBy(t *testing.T) {
	q := syncQB.SelectFrom("id, department", "table").GroupBy("department")
	q.Merge()
	if q.String() != "SELECT id, department FROM table GROUP BY department" {
		t.Error(testErrorText)
	}
}

func TestGroupByHavingCompare(t *testing.T) {
	q := syncQB.SelectFrom("id, department", "table").GroupBy("department").Having(qb.Compare("id", qb.GREATER, 5))
	q.Merge()
	if q.String() != "SELECT id, department FROM table GROUP BY department HAVING id > ?" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(q.Args(), []any{5}) {
		t.Error(testErrorArgs)
	}
}

func TestLimit(t *testing.T) {
	q := syncQB.SelectFrom("*", "table").Limit(5)
	q.Merge()
	if q.String() != "SELECT * FROM table LIMIT 5" {
		t.Error(testErrorText)
	}
}

func TestOffset(t *testing.T) {
	q := syncQB.SelectFrom("*", "table").Offset(5)
	q.Merge()
	if q.String() != "SELECT * FROM table OFFSET 5" {
		t.Error(testErrorText)
	}
}
