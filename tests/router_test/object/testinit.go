package object_test_1

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
)

type MockDatabase struct {
	db *sql.DB
}

func NewMockDatabase(db *sql.DB) *MockDatabase {
	return &MockDatabase{
		db: db,
	}
}

func (d *MockDatabase) SelectAll(tableName string) ([]map[string]interface{}, error) {
	rows, err := d.db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	res := []map[string]interface{}{}
	err = dbutils.ScanRows(rows, func(row map[string]interface{}) {
		res = append(res, row)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *MockDatabase) SelectWhereEqual(tableName string, colName string, val any) ([]map[string]interface{}, error) {
	rows, err := d.db.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, colName), val)
	if err != nil {
		return nil, err
	}
	res := []map[string]interface{}{}
	err = dbutils.ScanRows(rows, func(row map[string]interface{}) {
		res = append(res, row)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func initSelectWhereMock(mock sqlmock.Sqlmock, table string, val any, returnValues [][]driver.Value) {
	newRows := sqlmock.NewRows([]string{"id", "name", "ok"})
	for i := 0; i < len(returnValues); i++ {
		newRows.AddRows(returnValues[i])
	}
	mock.ExpectQuery(
		regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table)),
	).WithArgs(val).WillReturnRows(newRows)
}

func initSelectAllMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM table")).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "ok"}).
			AddRow(1, "TEST_NAME", true).
			AddRow(2, "TEST_NAME_1", true),
	)
}
