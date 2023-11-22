package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type Database struct {
	username string
	password string
	host     string
	port     string
	database string
	db       *sql.DB
	syncQ    interfaces.ISyncQueries
	asyncQ   interfaces.IAsyncQueries
}

func NewDatabase(username string, password string, host string, port string, database string) *Database {
	d := Database{username: username, password: password, host: host, port: port, database: database}
	d.SetSyncQueries(NewSyncQueries(&QueryBuild{}))
	d.SetAsyncQueries(NewAsyncQueries(&QueryBuild{}))
	return &d
}

// Connect connecting to a mysql database.
// Also, initialization of synchronous and asynchronous queries.
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

	d.syncQ.SetDB(db)
	d.asyncQ.SetSyncQueries(d.syncQ)
	return nil
}

// Close closes the connection to the database.
func (d *Database) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) SetSyncQueries(q interfaces.ISyncQueries) {
	d.syncQ = q
}

func (d *Database) SetAsyncQueries(q interfaces.IAsyncQueries) {
	d.asyncQ = q
}

func (d *Database) SyncQ() interfaces.ISyncQueries {
	return d.syncQ
}

func (d *Database) AsyncQ() interfaces.IAsyncQueries {
	return d.asyncQ
}

func (d *Database) DatabaseName() string {
	return d.database
}
