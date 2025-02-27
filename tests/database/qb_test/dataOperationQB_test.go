package qbtest

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	q := syncQB().Select("'Hello, world!'")
	q.Merge()
	if q.String() != "SELECT 'Hello, world!'" {
		t.Error(testErrorText)
	}
}

func TestSelectAs(t *testing.T) {
	q := syncQB().Select("'Hello, world!'").As("test")
	q.Merge()
	if q.String() != "SELECT 'Hello, world!' AS test" {
		t.Error(testErrorText)
	}
}

func TestSelectFrom(t *testing.T) {
	q := syncQB().SelectFrom("id", "tableName")
	q.Merge()
	if q.String() != "SELECT id FROM tableName" {
		t.Error(testErrorText)
	}
}

func TestInsert(t *testing.T) {
	q := syncQB().Insert("tableName", map[string]any{"id": 1, "name": "test"})
	q.Merge()
	if q.String() != "INSERT INTO tableName ( id, name ) VALUES ( ?, ? )" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(q.Args(), []any{1, "test"}) {
		t.Error(testErrorArgs)
	}
}

func TestUpdate(t *testing.T) {
	q := syncQB().Update("tableName", map[string]any{"id": 1, "name": "test"})
	q.Merge()
	if q.String() != "UPDATE tableName SET id = ?, name = ?" {
		t.Error(testErrorText)
	}
	if !reflect.DeepEqual(q.Args(), []any{1, "test"}) {
		t.Error(testErrorArgs)
	}
}

func TestDelete(t *testing.T) {
	q := syncQB().Delete("tableName")
	q.Merge()
	if q.String() != "DELETE FROM tableName" {
		t.Error(testErrorText)
	}
}
