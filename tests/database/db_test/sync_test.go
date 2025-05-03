package dbtest

import (
	"errors"
	"fmt"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

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
	err = db.Connect()
	if err != nil {
		t.Error(err)
	}
}

func TestSyncQuery(t *testing.T) {
	query, err := db.SyncQ().Query("SELECT `col1`, `col2`, `col3` FROM `dbtest` LIMIT 1")
	if err != nil {
		t.Error(err)
	}
	s := fmt.Sprintf("%v", query)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:111.22]]" {
		t.Errorf("The row values in the database of the expected row do not match.")
	}
}

func TestFUNCDatabaseResultIsEmpty_NotRaise(t *testing.T) {
	clearSyncTest()
	createSyncTest()
	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1='test1'")
	if err != nil {
		t.Error(err)
	}
	if err := dbutils.DatabaseResultNotEmpty(res); err != nil {
		t.Error(err)
	}
}

func TestFUNCDatabaseResultIsEmpty_Raise(t *testing.T) {
	clearSyncTest()
	createSyncTest()
	res, err := db.SyncQ().Query("SELECT * FROM dbtest WHERE col1='fff'")
	if err != nil {
		t.Error(err)
	}
	if err := dbutils.DatabaseResultNotEmpty(res); !errors.Is(err, dbutils.ErrDatabaseResultIsEmpty{}) {
		t.Error("error ErrDatabaseResultIsEmpty not raised")
	}
}
