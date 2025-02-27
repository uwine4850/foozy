package qbtest

import (
	"testing"

	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
)

func TestUnion(t *testing.T) {
	q := syncQB().SelectFrom("id", "table")
	q1 := syncQB().SelectFrom("id1", "table")
	resQ := qb.Union(q, q1)
	resQ.Merge()
	if resQ.String() != "(SELECT id FROM table) UNION (SELECT id1 FROM table)" {
		t.Error(testErrorText)
	}
}

func TestUnionAll(t *testing.T) {
	q := syncQB().SelectFrom("id", "table")
	q1 := syncQB().SelectFrom("id1", "table")
	resQ := qb.UnionAll(q, q1)
	resQ.Merge()
	if resQ.String() != "(SELECT id FROM table) UNION ALL (SELECT id1 FROM table)" {
		t.Error(testErrorText)
	}
}

func TestIntersect(t *testing.T) {
	q := syncQB().SelectFrom("id", "table")
	q1 := syncQB().SelectFrom("id1", "table")
	resQ := qb.Intersect(q, q1)
	resQ.Merge()
	if resQ.String() != "(SELECT id FROM table) INTERSECT (SELECT id1 FROM table)" {
		t.Error(testErrorText)
	}
}

func TestExcept(t *testing.T) {
	q := syncQB().SelectFrom("id", "table")
	q1 := syncQB().SelectFrom("id1", "table")
	resQ := qb.Except(q, q1)
	resQ.Merge()
	if resQ.String() != "(SELECT id FROM table) EXCEPT (SELECT id1 FROM table)" {
		t.Error(testErrorText)
	}
}
