package dbtest

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

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

func createDbAsyncTest() {
	_, err := db.SyncQ().Query("INSERT INTO `db_async_test` (`col1`, `col2`, `col3`) VALUES (?, ?, ?)",
		"test1", "2023-11-15", 111.22)
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Query("INSERT INTO `db_async_test` (`col1`, `col2`, `col3`) VALUES (?, ?, ?)",
		"test2", "2023-11-20", 222.11)
	if err != nil {
		panic(err)
	}
}

func clearDbAsyncTest() {
	_, err := db.SyncQ().Query("DELETE FROM db_async_test")
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
	clearDbAsyncTest()
	createDbAsyncTest()
	exitCode := m.Run()
	err = db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

type DbTestTable struct {
	Col1 string  `db:"col1"`
	Col2 string  `db:"col2"`
	Col3 float32 `db:"col3"`
}

func TestConnect(t *testing.T) {
	err := db.Connect()
	if err != nil {
		t.Error(err)
	}
}

func TestConnectErrorAndClose(t *testing.T) {
	err := db.Close()
	if err != nil {
		t.Error(err)
	}
	if err = db.Ping(); err != nil {
		connErr := database.ErrConnectionNotOpen{}
		if !errors.Is(err, connErr) {
			t.Errorf("The connection is open.")
		}
	}
}

func TestSyncQuery(t *testing.T) {
	db.Connect()
	query, err := db.SyncQ().Query("SELECT `col1`, `col2`, `col3` FROM `dbtest` LIMIT 1")
	if err != nil {
		t.Error(err)
	}
	s := fmt.Sprintf("%v", query)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:[49 49 49 46 50 50]]]" {
		t.Errorf("The row values in the database of the expected row do not match.")
	}
}

func TestSyncSelect(t *testing.T) {
	db.Connect()
	res1, err := db.SyncQ().Select([]string{"col1", "col2", "col3"}, "dbtest", dbutils.WHOutput{}, 0)
	if err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%v", res1)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:[49 49 49 46 50 50]] "+
		"map[col1:[116 101 115 116 50] col2:[50 48 50 51 45 49 49 45 50 48] col3:[50 50 50 46 49 49]]]" {
		t.Errorf("The result of sampling all fields is not the same as expected.")
	}
	res2, err := db.SyncQ().Select([]string{"col1", "col2", "col3"}, "dbtest", dbutils.WHOutput{}, 1)
	if err != nil {
		t.Error(err)
	}
	s2 := fmt.Sprintf("%v", res2)
	if s2 != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:[49 49 49 46 50 50]]]" {
		t.Errorf("The result of sampling all fields with a limit does not match the expected result.")
	}
}

func TestSyncSelectEquals(t *testing.T) {
	res3, err := db.SyncQ().Select([]string{"col1", "col2", "col3"}, "dbtest", dbutils.WHEquals(map[string]interface{}{
		"col1": "test2",
	}, "AND"), 1)
	if err != nil {
		t.Error(err)
	}
	s3 := fmt.Sprintf("%v", res3)
	if s3 != "[map[col1:[116 101 115 116 50] col2:[50 48 50 51 45 49 49 45 50 48] col3:222.11]]" {
		t.Errorf("The result of sampling fields c dbutils.WHEquals does not match the expected result.")
	}
}

func TestSyncSelectNotEquals(t *testing.T) {
	res4, err := db.SyncQ().Select([]string{"col1", "col2", "col3"}, "dbtest", dbutils.WHNotEquals(map[string]interface{}{
		"col1": "test2",
	}, "AND"), 1)
	if err != nil {
		panic(err)
	}
	s4 := fmt.Sprintf("%v", res4)
	if s4 != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:111.22]]" {
		t.Errorf("The result of c dbutils.WHNotEquals field sampling does not match the expected result.")
	}
}

func TestSyncSelectInSlice(t *testing.T) {
	res5, err := db.SyncQ().Select([]string{"col1", "col2", "col3"}, "dbtest", dbutils.WHInSlice(map[string][]interface{}{
		"col1": {"test1", "test2"},
	}, "AND"), 0)
	if err != nil {
		panic(err)
	}
	s5 := fmt.Sprintf("%v", res5)
	if s5 != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:111.22] "+
		"map[col1:[116 101 115 116 50] col2:[50 48 50 51 45 49 49 45 50 48] col3:222.11]]" {
		t.Errorf("The result of c dbutils.WHInSlice field sampling does not match the expected result.")
	}
}

func TestSyncSelectNotInSlice(t *testing.T) {
	res6, err := db.SyncQ().Select([]string{"col1", "col2", "col3"}, "dbtest", dbutils.WHNotInSlice(map[string][]interface{}{
		"col1": {"test1", "test2"},
	}, "AND"), 0)
	if err != nil {
		panic(err)
	}
	if res6 != nil {
		t.Errorf("The result of c dbutils.WHNotInSlice field sampling does not match the expected result.")
	}
}

func TestSyncInsert(t *testing.T) {
	clearDbTest()
	createDbTest()
	_, err := db.SyncQ().Insert("dbtest", map[string]interface{}{"col1": "text3", "col2": "2023-10-20", "col3": 10.22})
	if err != nil {
		panic(err)
	}
	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'text3'")
	if err != nil {
		panic(err)
	}
	if res == nil {
		t.Errorf("The Insert command failed.")
	}
}

func TestSyncInsertWithStruct(t *testing.T) {
	clearDbTest()
	createDbTest()
	insertStruct := DbTestTable{
		Col1: "ins1",
		Col2: "",
		Col3: 123,
	}
	insertStrcutValue, err := dbutils.ParamsValueFromStruct(&insertStruct, []string{"col2"})
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Insert("dbtest", insertStrcutValue)
	if err != nil {
		panic(err)
	}
	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'ins1'")
	if err != nil {
		panic(err)
	}
	if res == nil {
		t.Errorf("The Insert command failed.")
	}
}

func TestSyncCount(t *testing.T) {
	clearDbTest()
	createDbTest()
	count, err := db.SyncQ().Count([]string{"*"}, "dbtest", dbutils.WHOutput{}, 0)
	if err != nil {
		t.Error(err)
	}
	parseInt, err := dbutils.ParseInt(count[0]["COUNT(*)"])
	if err != nil {
		t.Error(err)
	}
	if parseInt < 2 {
		t.Errorf("The result of the command is not the same as expected.")
	}
}

func TestSyncDelete(t *testing.T) {
	clearDbTest()
	createDbTest()
	_, err := db.SyncQ().Delete("dbtest", dbutils.WHEquals(map[string]interface{}{
		"col1": "test1",
	}, "AND"))
	if err != nil {
		t.Error(err)
	}
	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'test1'")
	if err != nil {
		t.Error(err)
	}
	if res != nil {
		t.Errorf("The line has not been deleted.")
	}
}

func TestSyncUpdate(t *testing.T) {
	clearDbTest()
	createDbTest()
	_, err := db.SyncQ().Update("dbtest", map[string]any{"col1": "upd1", "col2": "2023-10-15", "col3": 1.1},
		dbutils.WHEquals(map[string]interface{}{"col1": "test2"}, "AND"))
	if err != nil {
		t.Error(err)
	}
	res, err := db.SyncQ().Select([]string{"*"}, "dbtest", dbutils.WHEquals(map[string]interface{}{
		"col1": "upd1", "col2": "2023-10-15",
	}, "AND"), 0)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Errorf("The row has not been updated.")
	}
}

func TestSyncUpdateWithStruct(t *testing.T) {
	clearDbTest()
	createDbTest()
	updateStruct := DbTestTable{
		Col1: "updStruct",
		Col2: "2023-10-16",
		Col3: 123,
	}
	params, err := dbutils.ParamsValueFromStruct(&updateStruct, []string{})
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Update("dbtest", params, dbutils.WHEquals(dbutils.WHValue{"col1": "test1"}, "AND"))
	if err != nil {
		t.Error(err)
	}
	res, err := db.SyncQ().Select([]string{"*"}, "dbtest", dbutils.WHEquals(map[string]interface{}{
		"col1": "updStruct", "col2": "2023-10-16",
	}, "AND"), 0)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Errorf("The field could not be updated using the structure.")
	}
}

func TestIncrement(t *testing.T) {
	clearDbTest()
	createDbTest()
	_, err := db.SyncQ().Increment("col4", "dbtest", dbutils.WHOutput{})
	if err != nil {
		panic(err)
	}
}

func TestSyncCommitTransaction(t *testing.T) {
	clearDbTest()
	createDbTest()
	db.BeginTransaction()
	_, err := db.SyncQ().Insert("dbtest", map[string]interface{}{"col1": "textComm", "col2": "2023-11-21", "col3": 10.24})
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Insert("dbtest", map[string]interface{}{"col1": "textComm1", "col2": "2023-11-21", "col3": 10.24})
	if err != nil {
		panic(err)
	}
	if err := db.CommitTransaction(); err != nil {
		panic(err)
	}

	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'textComm'")
	if err != nil {
		panic(err)
	}
	res1, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'textComm1'")
	if err != nil {
		panic(err)
	}
	if res == nil || res1 == nil {
		t.Errorf("The commit transaction failed.")
	}
}

func TestSyncRollbackTransaction(t *testing.T) {
	clearDbTest()
	createDbTest()
	db.BeginTransaction()
	_, err := db.SyncQ().Insert("dbtest", map[string]interface{}{"col1": "textBack", "col2": "2023-11-21", "col3": 10.24})
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Insert("dbtest", map[string]interface{}{"col1": "textBack1", "col2": "2023-11-21", "col3": 10.24})
	if err != nil {
		panic(err)
	}
	if err := db.RollBackTransaction(); err != nil {
		panic(err)
	}

	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'textBack'")
	if err != nil {
		panic(err)
	}
	res1, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1 = 'textBack1'")
	if err != nil {
		panic(err)
	}
	if res != nil || res1 != nil {
		t.Errorf("The rollback transaction failed.")
	}
}
