package dbtest

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"testing"
)

func TestAsyncSelect(t *testing.T) {
	db.AsyncQ().AsyncSelect("s", []string{"col1", "col2", "col3"}, "db_async_test", dbutils.WHOutput{}, 0)
	db.AsyncQ().AsyncSelect("s1", []string{"col1", "col2", "col3"}, "db_async_test", dbutils.WHOutput{}, 1)
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	res1, _ := db.AsyncQ().LoadAsyncRes("s1")
	s := fmt.Sprintf("%v", res.Res)
	if res.Error != nil {
		t.Error(res.Error)
	}
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:[49 49 49 46 50 50]] "+
		"map[col1:[116 101 115 116 50] col2:[50 48 50 51 45 49 49 45 50 48] col3:[50 50 50 46 49 49]]]" {
		t.Errorf("The result of sampling all fields is not the same as expected.")
	}
	s2 := fmt.Sprintf("%v", res1.Res)
	if res1.Error != nil {
		t.Error(res1.Error)
	}
	if s2 != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:[49 49 49 46 50 50]]]" {
		t.Errorf("The result of sampling all fields with a limit does not match the expected result.")
	}
}

func TestAsyncQuery(t *testing.T) {
	db.AsyncQ().AsyncQuery("s", "SELECT `col1`, `col2`, `col3` FROM `db_async_test` LIMIT 1")
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	s := fmt.Sprintf("%v", res.Res)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:[49 49 49 46 50 50]]]" {
		t.Errorf("The row values in the database of the expected row do not match.")
	}
}

func TestAsyncSelectEquals(t *testing.T) {
	db.AsyncQ().AsyncSelect("s", []string{"col1", "col2", "col3"}, "db_async_test", dbutils.WHEquals(map[string]interface{}{
		"col1": "test2",
	}, "AND"), 1)
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	s := fmt.Sprintf("%v", res.Res)
	if s != "[map[col1:[116 101 115 116 50] col2:[50 48 50 51 45 49 49 45 50 48] col3:222.11]]" {
		t.Errorf("The result of sampling fields c dbutils.WHEquals does not match the expected result.")
	}
}

func TestAsyncSelectNotEquals(t *testing.T) {
	db.AsyncQ().AsyncSelect("s", []string{"col1", "col2", "col3"}, "db_async_test", dbutils.WHNotEquals(map[string]interface{}{
		"col1": "test2",
	}, "AND"), 1)
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	s := fmt.Sprintf("%v", res.Res)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:111.22]]" {
		t.Errorf("The result of c dbutils.WHNotEquals field sampling does not match the expected result.")
	}
}

func TestAsyncSelectInSlice(t *testing.T) {
	db.AsyncQ().AsyncSelect("s", []string{"col1", "col2", "col3"}, "db_async_test", dbutils.WHInSlice(map[string][]interface{}{
		"col1": {"test1", "test2"},
	}, "AND"), 0)
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	s := fmt.Sprintf("%v", res.Res)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:111.22] "+
		"map[col1:[116 101 115 116 50] col2:[50 48 50 51 45 49 49 45 50 48] col3:222.11]]" {
		t.Errorf("The result of c dbutils.WHInSlice field sampling does not match the expected result.")
	}
}

func TestAsyncSelectNotInSlice(t *testing.T) {
	db.AsyncQ().AsyncSelect("s", []string{"col1", "col2", "col3"}, "db_async_test", dbutils.WHNotInSlice(map[string][]interface{}{
		"col1": {"test1", "test2"},
	}, "AND"), 0)
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	if res.Res != nil {
		t.Errorf("The result of c dbutils.WHNotInSlice field sampling does not match the expected result.")
	}
}

func TestAsyncInsert(t *testing.T) {
	db.AsyncQ().AsyncInsert("s", "db_async_test", map[string]interface{}{"col1": "text3", "col2": "2023-10-20", "col3": 10.22})
	db.AsyncQ().Wait()
	db.AsyncQ().AsyncQuery("s1", "SELECT * FROM db_async_test WHERE col1 = 'text3'")
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	res1, _ := db.AsyncQ().LoadAsyncRes("s1")
	if res1.Error != nil {
		t.Error(res1.Error)
	}
	if res1.Res == nil {
		t.Errorf("The Insert command failed.")
	}
}

func TestAsyncCount(t *testing.T) {
	db.AsyncQ().AsyncCount("s", []string{"*"}, "db_async_test", dbutils.WHOutput{}, 0)
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	parseInt, err := dbutils.ParseInt(res.Res[0]["COUNT(*)"])
	if err != nil {
		t.Error(err)
	}
	if parseInt < 2 {
		t.Errorf("The result of the command is not the same as expected.")
	}
}

func TestAsyncDelete(t *testing.T) {
	db.AsyncQ().AsyncDelete("s", "db_async_test", dbutils.WHEquals(map[string]interface{}{
		"col1": "test1",
	}, "AND"))
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	db.AsyncQ().AsyncQuery("s1", "SELECT * FROM db_async_test WHERE col1 = 'test1'")
	db.AsyncQ().Wait()
	res1, _ := db.AsyncQ().LoadAsyncRes("s1")
	if res1.Error != nil {
		t.Error(res1.Error)
	}
	if res1.Res != nil {
		t.Errorf("The line has not been deleted.")
	}
}

func TestAsyncUpdate(t *testing.T) {
	db.AsyncQ().AsyncUpdate("s", "db_async_test", []dbutils.DbEquals{
		{Name: "col1", Value: "upd1"},
		{Name: "col2", Value: "2023-10-15"},
		{Name: "col3", Value: 1.1},
	}, dbutils.WHEquals(map[string]interface{}{"col1": "test2"}, "AND"))
	db.AsyncQ().Wait()
	res, _ := db.AsyncQ().LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	db.AsyncQ().AsyncSelect("s1", []string{"*"}, "db_async_test", dbutils.WHEquals(map[string]interface{}{
		"col1": "upd1", "col2": "2023-10-15",
	}, "AND"), 0)
	db.AsyncQ().Wait()
	res1, _ := db.AsyncQ().LoadAsyncRes("s1")
	if res1.Error != nil {
		t.Error(res1.Error)
	}
	if res1.Res == nil {
		t.Errorf("The row has not been updated.")
	}
}
