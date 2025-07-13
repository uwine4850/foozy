package database_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/database"
	databasemock "github.com/uwine4850/foozy/tests1/database_mock"
)

var dbmock *databasemock.MysqlDatabase

func TestMain(m *testing.M) {
	syncQ := database.NewSyncQueries()
	asyncQ := database.NewAsyncQueries(syncQ)
	dbmock = databasemock.NewMysqlDatabase(syncQ, asyncQ)
	if err := dbmock.Open(); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSyncQQuery(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "NAME")
	dbmock.Mock().ExpectQuery(regexp.QuoteMeta("SELECT * FROM table")).
		WillReturnRows(rows)
	_, err := dbmock.SyncQ().Query("SELECT * FROM table")
	if err != nil {
		t.Error(err)
	}
}

func TestSyncQExec(t *testing.T) {
	dbmock.Mock().ExpectExec(regexp.QuoteMeta("INSERT INTO table (id) VALUE (?)")).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := dbmock.SyncQ().Exec("INSERT INTO table (id) VALUE (?)", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestSyncQTransaction(t *testing.T) {
	dbmock.Mock().ExpectBegin()
	dbmock.Mock().ExpectExec(regexp.QuoteMeta("INSERT INTO table (id) VALUE (?)")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	dbmock.Mock().ExpectCommit()
	transaction, err := dbmock.NewTransaction()
	if err != nil {
		t.Error(err)
	}
	if err := transaction.BeginTransaction(); err != nil {
		t.Error(err)
	}
	_, err = transaction.SyncQ().Exec("INSERT INTO table (id) VALUE (?)", 1)
	if err != nil {
		t.Error(err)
	}
	if err := transaction.CommitTransaction(); err != nil {
		t.Error(err)
	}
}

func TestAsyncQQuery(t *testing.T) {
	dbmock.Mock().MatchExpectationsInOrder(false)

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)
	dbmock.Mock().ExpectQuery(regexp.QuoteMeta("SELECT id FROM table")).
		WillReturnRows(rows)
	rows2 := sqlmock.NewRows([]string{"name"})
	rows2.AddRow("NAME")
	dbmock.Mock().ExpectQuery(regexp.QuoteMeta("SELECT name FROM table")).
		WillReturnRows(rows2)

	asyncQ, err := dbmock.NewAsyncQ()
	if err != nil {
		t.Error(err)
	}
	asyncQ.Query("q1", "SELECT id FROM table")
	asyncQ.Query("q2", "SELECT name FROM table")
	asyncQ.Wait()
	resQ1, _ := asyncQ.LoadAsyncRes("q1")
	if resQ1.Error != nil {
		t.Error(resQ1.Error)
	}
	if resQ1.Res[0]["id"].(int64) != 1 {
		t.Error("the result of the first query does not match the expectation")
	}
	resQ2, _ := asyncQ.LoadAsyncRes("q2")
	if resQ2.Error != nil {
		t.Error(resQ2.Error)
	}
	if resQ2.Res[0]["name"].(string) != "NAME" {
		t.Error("the result of the second query does not match the expectation")
	}
}

func TestAsyncQExec(t *testing.T) {
	dbmock.Mock().MatchExpectationsInOrder(false)

	dbmock.Mock().ExpectExec(regexp.QuoteMeta("INSERT INTO table (id) VALUE (?)")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	dbmock.Mock().ExpectExec(regexp.QuoteMeta("INSERT INTO table (name) VALUE (?)")).
		WithArgs("NAME").
		WillReturnResult(sqlmock.NewResult(1, 1))

	asyncQ, err := dbmock.NewAsyncQ()
	if err != nil {
		t.Error(err)
	}
	asyncQ.Exec("q1", "INSERT INTO table (id) VALUE (?)", 1)
	asyncQ.Exec("q2", "INSERT INTO table (name) VALUE (?)", "NAME")
	asyncQ.Wait()
	resQ1, _ := asyncQ.LoadAsyncRes("q1")
	if resQ1.Error != nil {
		t.Error(resQ1.Error)
	}
	resQ2, _ := asyncQ.LoadAsyncRes("q2")
	if resQ2.Error != nil {
		t.Error(resQ2.Error)
	}
}

func TestAsyncQTransaction(t *testing.T) {
	dbmock.Mock().MatchExpectationsInOrder(false)

	dbmock.Mock().ExpectBegin()
	dbmock.Mock().ExpectExec(regexp.QuoteMeta("INSERT INTO table (id) VALUE (?)")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	dbmock.Mock().ExpectExec(regexp.QuoteMeta("INSERT INTO table (name) VALUE (?)")).
		WithArgs("NAME").
		WillReturnResult(sqlmock.NewResult(1, 1))
	dbmock.Mock().ExpectCommit()

	transaction, err := dbmock.NewTransaction()
	if err != nil {
		t.Error(err)
	}
	if err := transaction.BeginTransaction(); err != nil {
		t.Error(err)
	}
	asyncTransaction, err := transaction.NewAsyncQ()
	if err != nil {
		t.Error(err)
	}
	asyncTransaction.Exec("q1", "INSERT INTO table (id) VALUE (?)", 1)
	asyncTransaction.Exec("q2", "INSERT INTO table (name) VALUE (?)", "NAME")
	asyncTransaction.Wait()
	if err := transaction.CommitTransaction(); err != nil {
		t.Error(err)
	}
	resQ1, _ := asyncTransaction.LoadAsyncRes("q1")
	if resQ1.Error != nil {
		t.Error(resQ1.Error)
	}
	resQ2, _ := asyncTransaction.LoadAsyncRes("q2")
	if resQ2.Error != nil {
		t.Error(resQ2.Error)
	}
}
