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
	syncQ  interfaces.SyncQ
	asyncQ interfaces.AsyncQ
	mock   sqlmock.Sqlmock
}

func NewMysqlDatabase(syncQ interfaces.SyncQ, asyncQ interfaces.AsyncQ) *MysqlDatabase {
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

func (d *MysqlDatabase) NewTransaction() (interfaces.DatabaseTransaction, error) {
	return database.NewMysqlTransaction(d.db, d.syncQ, d.asyncQ)
}

func (d *MysqlDatabase) SyncQ() interfaces.SyncQ {
	return d.syncQ
}

func (d *MysqlDatabase) NewAsyncQ() (interfaces.AsyncQ, error) {
	aq, err := d.asyncQ.New()
	if err != nil {
		return nil, err
	}
	return aq.(interfaces.AsyncQ), nil
}

func (d *MysqlDatabase) Mock() sqlmock.Sqlmock {
	return d.mock
}
