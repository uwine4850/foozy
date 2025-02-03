package dbtest

import (
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/tests/common/tconf"
)

var db = database.NewDatabase(tconf.DbArgs)

func createSyncTest() {
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

func clearSyncTest() {
	_, err := db.SyncQ().Query("DELETE FROM dbtest")
	if err != nil {
		panic(err)
	}
}

func createAsyncTest() {
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

func clearAsyncTest() {
	_, err := db.SyncQ().Query("DELETE FROM db_async_test")
	if err != nil {
		panic(err)
	}
}

func create() {
	createSyncTest()
	createAsyncTest()
}

func clear() {
	clearSyncTest()
	clearAsyncTest()
}

func TestMain(m *testing.M) {
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	clear()
	create()
	exitCode := m.Run()
	err = db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}
