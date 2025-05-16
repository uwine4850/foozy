package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uwine4850/foozy/pkg/config"
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
// For each transaction a new instance of the [interfaces.ITransaction] object is created,
// so each transaction is executed in its own scope and is completely safe.
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

func NewDatabase(args DbArgs, syncQ interfaces.ISyncQueries, asyncQ interfaces.IAsyncQueries) *Database {
	d := Database{username: args.Username, password: args.Password, host: args.Host, port: args.Port, database: args.DatabaseName}
	d.syncQ = syncQ
	d.asyncQ = asyncQ
	return &d
}

// Open connecting to a mysql database.
// Also, initialization of synchronous and asynchronous queries.
func (d *Database) Open() error {
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
	d.syncQ.SetDB(&DbQuery{DB: db})
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
			return err
		}
		d.tx = nil
	}
	return nil
}

// NewTransaction creates a new transaction instance.
func (d *Database) NewTransaction() interfaces.ITransaction {
	return NewMysqlTransaction(d.db, d.syncQ, d.asyncQ)
}

// SyncQ getting access to synchronous requests.
func (d *Database) SyncQ() interfaces.ISyncQueries {
	return d.syncQ
}

func (d *Database) NewAsyncQ() (interfaces.IAsyncQueries, error) {
	aq, err := d.asyncQ.New()
	if err != nil {
		return nil, err
	}
	return aq.(interfaces.IAsyncQueries), nil
}

// MysqlTransaction An object that performs transactions
// to the mysql database.
// This object is used only for one transaction, for each
// next transaction a new instance of the object must be created.
type MysqlTransaction struct {
	db     *sql.DB
	tx     *sql.Tx
	syncQ  interfaces.ISyncQueries
	asyncQ interfaces.IAsyncQueries
}

// NewMysqlTransaction creates a new [MysqlTransaction] escamp.
// For correct creation, it is necessary to pass an already open
// connection to the database.
func NewMysqlTransaction(db *sql.DB, syncQ interfaces.ISyncQueries, asyncQ interfaces.IAsyncQueries) *MysqlTransaction {
	return &MysqlTransaction{
		db:     db,
		syncQ:  syncQ,
		asyncQ: asyncQ,
	}
}

// BeginTransaction starts the transaction.
// Only one transaction can be started per object instance.
func (t *MysqlTransaction) BeginTransaction() error {
	if t.tx != nil {
		return errors.New("transaction already started")
	}
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	t.syncQ.SetDB(&DbTxQuery{Tx: tx})
	t.asyncQ.SetSyncQueries(t.syncQ)
	t.tx = tx
	return nil
}

// CommitTransaction writes the transaction to the database.
func (t *MysqlTransaction) CommitTransaction() error {
	if err := t.tx.Commit(); err != nil {
		return err
	}
	t.tx = nil
	return nil
}

// RollBackTransaction undoes any changes that were made
// during the transaction.
// That is, after executing the [BeginTransaction] method.
func (t *MysqlTransaction) RollBackTransaction() error {
	if err := t.tx.Rollback(); err != nil {
		return err
	}
	t.tx = nil
	return nil
}

// SyncQ getting access to synchronous requests.
func (t *MysqlTransaction) SyncQ() interfaces.ISyncQueries {
	return t.syncQ
}

// AsyncQ getting access to asynchronous requests.
func (t *MysqlTransaction) NewAsyncQ() (interfaces.IAsyncQueries, error) {
	aq, err := t.asyncQ.New()
	if err != nil {
		return nil, err
	}
	return aq.(interfaces.IAsyncQueries), nil
}

// DbQuery standard database queries. They are used *sql.DB.
// Requests are executed as usual.
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
	return map[string]interface{}{"insertID": id, "rowsAffected": rowsId}, nil
}

// DbTxQuery queries that can be rolled back. Used *sql.Tx.
// This object will perform queries with the [*sql.Tx]
// object that is used for transactions.
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
	return map[string]interface{}{"id": id, "rows": rowsId}, nil
}

// InitDatabasePool initializes a database pool.
// Only one pool is created, which is specified in the
// [Default.Database.MainConnectionPoolName] settings.
// Once created, the pool is locked. Therefore, you need to
// initialize it manually if you need more connections.
func InitDatabasePool(manager interfaces.IManager, db interfaces.IDatabase) error {
	name := config.LoadedConfig().Default.Database.MainConnectionPoolName
	if err := manager.Database().AddConnection(name, db); err != nil {
		return err
	}
	manager.Database().Lock()
	return nil
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
