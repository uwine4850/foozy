package databasemock

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type MysqlDatabase struct {
	db     *sql.DB
	tx     *sql.Tx
	syncQ  interfaces.ISyncQueries
	asyncQ interfaces.IAsyncQueries
	mock   sqlmock.Sqlmock
}

func NewMysqlDatabase(syncQ interfaces.ISyncQueries, asyncQ interfaces.IAsyncQueries) *MysqlDatabase {
	return &MysqlDatabase{
		syncQ:  syncQ,
		asyncQ: asyncQ,
	}
}

func (d *MysqlDatabase) Open() error {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		return err
	}
	d.db = db
	d.mock = sqlMock
	d.syncQ.SetDB(&database.DbQuery{DB: db})
	return nil
}

func (d *MysqlDatabase) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}
	if d.tx != nil {
		if err := d.tx.Rollback(); err != nil {
			return err
		}
		d.tx = nil
	}
	return nil
}

func (d *MysqlDatabase) NewTransaction() (interfaces.ITransaction, error) {
	return database.NewMysqlTransaction(d.db, d.syncQ, d.asyncQ)
}

func (d *MysqlDatabase) SyncQ() interfaces.ISyncQueries {
	return d.syncQ
}

func (d *MysqlDatabase) NewAsyncQ() (interfaces.IAsyncQueries, error) {
	aq, err := d.asyncQ.New()
	if err != nil {
		return nil, err
	}
	return aq.(interfaces.IAsyncQueries), nil
}

func (d *MysqlDatabase) Mock() sqlmock.Sqlmock {
	return d.mock
}
