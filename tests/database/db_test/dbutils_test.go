package dbtest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

func TestParseString(t *testing.T) {
	clearSyncTest()
	createSyncTest()
	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
	if err != nil {
		panic(err)
	}
	p := dbutils.ParseString(res[0]["col1"])
	if p != "test1" {
		t.Errorf("The result of the method is not the same as expected.")
	}
}

func TestParseInt(t *testing.T) {
	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	parseInt, err := dbutils.ParseInt(res[0]["id"])
	if err != nil {
		panic(err)
	}
	if parseInt <= 0 {
		t.Errorf("Error executing the dbutils.ParseInt method.")
	}
}

func TestParseDateTime(t *testing.T) {
	clearSyncTest()
	createSyncTest()
	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
	if err != nil {
		panic(err)
	}
	time, err := dbutils.ParseDateTime("2006-01-02", res[0]["col2"])
	if err != nil {
		panic(err)
	}
	if time.String() == "2023-11-15" {
		t.Errorf("Error executing dbutils.ParseDateTime method.")
	}
}

func TestParseFloat(t *testing.T) {
	clearSyncTest()
	createSyncTest()
	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
	if err != nil {
		panic(err)
	}
	float, err := dbutils.ParseFloat(res[0]["col3"])
	if err != nil {
		panic(err)
	}
	if float != 111.22 {
		t.Errorf("Execution error of dbutils.ParseFloat method.")
	}
}

func TestParseEquals(t *testing.T) {
	var equals []dbutils.DbEquals
	equals = append(equals, dbutils.DbEquals{
		Name:  "e1",
		Value: 1,
	})
	equals = append(equals, dbutils.DbEquals{
		Name:  "e2",
		Value: 2,
	})
	parseEquals, args := dbutils.ParseEquals(equals, "AND")
	if parseEquals != "e1 = ? AND e2 = ?" {
		t.Errorf("The timing doesn't match the expectation. Expected e1 = ? AND e2 = ?, received %s", parseEquals)
	}
	i := []interface{}{1, 2}
	if !reflect.DeepEqual(args, i) {
		t.Errorf("Arguments don't match the expectation.")
	}
}

func TestParseParams(t *testing.T) {
	params, args := dbutils.ParseParams(map[string]interface{}{"p1": 1, "p2": 2})
	if !fslice.SliceContains(params, "p1") || !fslice.SliceContains(params, "p2") {
		t.Errorf("The keys are not as expected.")
	}
	if !fslice.SliceContains(args, 1) || !fslice.SliceContains(args, 2) {
		t.Errorf("Arguments don't match the expectation.")
	}
}

func TestRepeatValues(t *testing.T) {
	values := dbutils.RepeatValues(3, ",")
	if values != "?, ?, ?" {
		t.Errorf("The result does not match the expectation.")
	}
}
