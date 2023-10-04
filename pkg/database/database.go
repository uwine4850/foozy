package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uwine4850/foozy/internal/interfaces"
)

type Database struct {
	username string
	password string
	host     string
	port     string
	database string
	db       *sql.DB
	SyncQ    interfaces.ISyncQueries
	AsyncQ   interfaces.IAsyncQueries
}

func NewDatabase(username string, password string, host string, port string, database string) *Database {
	d := Database{username: username, password: password, host: host, port: port, database: database}
	d.SetSyncQueries(&SyncQueries{})
	d.SetAsyncQueries(&AsyncQueries{})
	return &d
}

func (d *Database) Connect() error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.username, d.password, d.host, d.port, d.database)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	d.db = db

	d.SyncQ.SetDB(db)
	d.AsyncQ.SetSyncQueries(d.SyncQ)
	return nil
}

func (d *Database) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) SetSyncQueries(q interfaces.ISyncQueries) {
	d.SyncQ = q
}

func (d *Database) SetAsyncQueries(q interfaces.IAsyncQueries) {
	d.AsyncQ = q
}
