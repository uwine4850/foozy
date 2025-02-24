package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type DbArgs struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

// Database structure for accessing the database.
// It can send both synchronous and asynchronous queries.
// IMPORTANT: after the end of work it is necessary to close the connection using Close method.
//
// Principle of [db] and [tx] swapping:
// After the initial initialization, [db] is used. These are standard database queries. If the [BeginTransaction] method is used,
// the [db] instance will be replaced by the [tx] instance. Then the [CommitTransaction] method changes them back.
// These instances use the same interface, so they are directly used by the [ISyncQueries] interface, which in turn is used by
// the [IAsyncQueries] interface.
// The main difference between [db] and [tx] objects is that the latter is used for the ability to cancel database transactions.
type Database struct {
	username string
	password string
	host     string
	port     string
	database string
	db       *sql.DB
	tx       *sql.Tx
	syncQ    interfaces.ISyncQueries
	asyncQ   interfaces.IAsyncQueries
}

func NewDatabase(args DbArgs) *Database {
	d := Database{username: args.Username, password: args.Password, host: args.Host, port: args.Port, database: args.DatabaseName}
	d.SetSyncQueries(NewSyncQueries(&QueryBuild{}))
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
	// AsyncQueries should be installed exactly when a new connection to the database is made.
	// This is mandatory because AsyncQueries contains the results of queries,
	// so you need to create a new instance for this object to make the data unique for each user.
	// SyncQueries does not store query data in the object, so there is no need to initialize it every connection.
	d.SetAsyncQueries(NewAsyncQueries(&QueryBuild{}))
	d.syncQ.SetDB(&DbQuery{DB: db})
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
	if d.tx != nil {
		if err := d.tx.Rollback(); err != nil {
			return nil
		}
		d.tx = nil
	}
	return nil
}

// BeginTransaction begins executing a transaction.
// Changes the executable object for database queries, so all queries that run under this method will use *sql.Tx.
func (d *Database) BeginTransaction() {
	tx, err := d.db.Begin()
	if err != nil {
		panic(err)
	}
	d.syncQ.SetDB(&DbTxQuery{Tx: tx})
	d.asyncQ.SetSyncQueries(d.syncQ)
	d.tx = tx
}

// CommitTransaction records changes in the database.
// This method ends the transaction and changes the executable query object from *sql.Tx to *sql.DB.
// Therefore, all the following queries are executed using *sql.DB.
func (d *Database) CommitTransaction() error {
	if err := d.tx.Commit(); err != nil {
		return err
	}
	d.tx = nil
	d.syncQ.SetDB(&DbQuery{DB: d.db})
	d.asyncQ.SetSyncQueries(d.syncQ)
	return nil
}

// RollBackTransaction rolls back changes made by commands AFTER CALLING BeginTransaction().
// This method ends the transaction and changes the executable query object from *sql.Tx to *sql.DB.
func (d *Database) RollBackTransaction() error {
	if err := d.tx.Rollback(); err != nil {
		return err
	}
	d.tx = nil
	d.syncQ.SetDB(&DbQuery{DB: d.db})
	d.asyncQ.SetSyncQueries(d.syncQ)
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

// DbQuery standard database queries. They are used *sql.DB.
type DbQuery struct {
	DB *sql.DB
}

func (d *DbQuery) Query(query string, args ...any) ([]map[string]interface{}, error) {
	sqlRows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	rows, err := scanRows(sqlRows)
	if err != nil {
		return nil, err
	}

	err = sqlRows.Close()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DbQuery) Exec(query string, args ...any) (map[string]interface{}, error) {
	result, err := d.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	rowsId, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"id": id, "rows": rowsId}, nil
}

// DbTxQuery queries that can be rolled back. Used *sql.Tx.
type DbTxQuery struct {
	Tx *sql.Tx
}

func (d *DbTxQuery) Query(query string, args ...any) ([]map[string]interface{}, error) {
	sqlRows, err := d.Tx.Query(query, args...)
	if err != nil {
		return nil, err
	}

	rows, err := scanRows(sqlRows)
	if err != nil {
		return nil, err
	}

	err = sqlRows.Close()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DbTxQuery) Exec(query string, args ...any) (map[string]interface{}, error) {
	result, err := d.Tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	rowsId, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"id": &id, "rows": &rowsId}, nil
}

func scanRows(sqlRows *sql.Rows) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	if err := dbutils.ScanRows(sqlRows, func(row map[string]interface{}) {
		rows = append(rows, row)
	}); err != nil {
		return nil, err
	}
	return rows, nil
}

type ErrConnectionNotOpen struct {
}

func (receiver ErrConnectionNotOpen) Error() string {
	return "The connection is not open."
}
