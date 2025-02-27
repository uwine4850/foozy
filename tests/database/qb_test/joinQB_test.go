package qbtest

import (
	"testing"

	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
)

func TestInnerJoin(t *testing.T) {
	q := syncQB().SelectFrom("auth.id, user_roles.id", "auth").InnerJoin("user_roles", qb.NoArgsCompare("auth.id", qb.EQUAL, "user_roles.id"))
	q.Merge()
	if q.String() != "SELECT auth.id, user_roles.id FROM auth INNER JOIN user_roles ON auth.id = user_roles.id" {
		t.Error(testErrorText)
	}
}

func TestLeftJoin(t *testing.T) {
	q := syncQB().SelectFrom("auth.id, user_roles.id", "auth").LeftJoin("user_roles", qb.NoArgsCompare("auth.id", qb.EQUAL, "user_roles.id"))
	q.Merge()
	if q.String() != "SELECT auth.id, user_roles.id FROM auth LEFT JOIN user_roles ON auth.id = user_roles.id" {
		t.Error(testErrorText)
	}
}

func TestRightJoin(t *testing.T) {
	q := syncQB().SelectFrom("auth.id, user_roles.id", "auth").RightJoin("user_roles", qb.NoArgsCompare("auth.id", qb.EQUAL, "user_roles.id"))
	q.Merge()
	if q.String() != "SELECT auth.id, user_roles.id FROM auth RIGHT JOIN user_roles ON auth.id = user_roles.id" {
		t.Error(testErrorText)
	}
}
