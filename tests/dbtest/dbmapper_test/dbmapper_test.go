package dbmappertest_test

import (
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbmapper"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

var dbArgs = database.DbArgs{
	Username: "root", Password: "1111", Host: "localhost", Port: "3408", DatabaseName: "foozy_test",
}
var db = database.NewDatabase(dbArgs)

func createDbTest() {
	_, err := db.SyncQ().Query("INSERT INTO `dbtest` (`col1`, `col2`, `col3`) VALUES (?, ?, ?)",
		"test1", "2023-11-15", 111.22)
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Query("INSERT INTO `dbtest` (`col1`, `col2`, `col3`) VALUES (?, ?, ?)",
		"test2", "2023-11-20", 222.11)
	if err != nil {
		panic(err)
	}
}

func clearDbTest() {
	_, err := db.SyncQ().Query("DELETE FROM dbtest")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	clearDbTest()
	createDbTest()
	exitCode := m.Run()
	err = db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

type DbTestMapper struct {
	Col1 string `db:"col1"`
	Col2 string `db:"col2"`
	Col3 string `db:"col3"`
	Col4 string `db:"col4"`
}

func TestDbMapperUseStruct(t *testing.T) {
	clearDbTest()
	createDbTest()
	res, err := db.SyncQ().Select([]string{"*"}, "dbtest", dbutils.WHOutput{}, 0)
	if err != nil {
		t.Error(err)
	}
	var dbTestMapper []DbTestMapper
	mapper := dbmapper.NewMapper(res, &dbTestMapper)
	if err := mapper.Fill(); err != nil {
		t.Error(err)
	}
	if len(dbTestMapper) == 0 {
		t.Error("DbMapper.Output must not be empty")
	}
	map1 := DbTestMapper{Col1: "test1", Col2: "2023-11-15", Col3: "111.22"}
	if dbTestMapper[0] != map1 {
		t.Error("DbMapper.Output value does not match expected")
	}
	map2 := DbTestMapper{Col1: "test2", Col2: "2023-11-20", Col3: "222.11"}
	if dbTestMapper[1] != map2 {
		t.Error("DbMapper.Output value does not match expected")
	}
}

func TestDbMapperUseMap(t *testing.T) {
	clearDbTest()
	createDbTest()
	res, err := db.SyncQ().Select([]string{"*"}, "dbtest", dbutils.WHOutput{}, 0)
	if err != nil {
		t.Error(err)
	}
	var dbTestMapper = []map[string]string{}
	mapper := dbmapper.NewMapper(res, &dbTestMapper)
	if err := mapper.Fill(); err != nil {
		t.Error(err)
	}
	if len(dbTestMapper) == 0 {
		t.Error("DbMapper.Output must not be empty")
	}

	map1 := map[string]string{"col1": "test1", "col2": "2023-11-15", "col3": "111.22", "col4": "", "id": ""}
	if len(dbTestMapper[0]) != len(map1) {
		t.Error("DbMapper.Output map len does not match expected")
	}
	map2 := map[string]string{"col1": "test2", "col2": "2023-11-20", "col3": "222.11", "col4": "", "id": ""}
	if len(dbTestMapper[1]) != len(map2) {
		t.Error("DbMapper.Output map len does not match expected")
	}
	if dbTestMapper[0]["col1"] != map1["col1"] {
		t.Error("DbMapper.Output value does not match expected")
	}
}
