package utilstest

import (
	"testing"

	"github.com/uwine4850/foozy/pkg/utils/fmap"
)

func TestMergeMap(t *testing.T) {
	m1 := map[string]int{"1": 1}
	m2 := map[string]int{"2": 2}
	fmap.MergeMap(&m1, m2)
	if m1["1"] != 1 && m1["2"] != 2 {
		t.Error("The map does not match expectations.")
	}
}

func TestCompare(t *testing.T) {
	m1 := map[string]int{"1": 1, "2": 2, "exc": 111}
	m2 := map[string]int{"1": 1, "2": 2, "exc": 222}
	m3 := map[string]int{"1": 11, "2": 22}
	if !fmap.Compare(&m1, &m2, []string{"exc"}) {
		t.Error("Map comparison error.")
	}
	if fmap.Compare(&m1, &m3, nil) {
		t.Error("Map comparison did not return false.")
	}
}
