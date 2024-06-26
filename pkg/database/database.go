package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type DbArgs struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

type ErrConnectionNotOpen struct {
}

func (receiver ErrConnectionNotOpen) Error() string {
	return "The connection is not open."
}

// Database structure for accessing the database.
// It can send both synchronous and asynchronous queries.
// IMPORTANT: after the end of work it is necessary to close the connection using Close method.
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

func NewDatabase(args DbArgs) *Database {
	d := Database{username: args.Username, password: args.Password, host: args.Host, port: args.Port, database: args.DatabaseName}
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

func (d *Database) Ping() error {
	if d.db == nil {
		return ErrConnectionNotOpen{}
	}
	err := d.db.Ping()
	if err != nil {
		return ErrConnectionNotOpen{}
	}
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

// SetSyncQueries sets the synchronous query interface.
func (d *Database) SetSyncQueries(q interfaces.ISyncQueries) {
	d.syncQ = q
}

// SetAsyncQueries sets the asynchronous query interface.
func (d *Database) SetAsyncQueries(q interfaces.IAsyncQueries) {
	d.asyncQ = q
}

// SyncQ getting access to synchronous requests.
func (d *Database) SyncQ() interfaces.ISyncQueries {
	return d.syncQ
}

// AsyncQ getting access to asynchronous requests.
func (d *Database) AsyncQ() interfaces.IAsyncQueries {
	return d.asyncQ
}

// DatabaseName Getting the database name.
func (d *Database) DatabaseName() string {
	return d.database
}
