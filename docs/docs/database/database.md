## database
Implementation of Mysql database interfaces.

### MysqlDatabase
Implementation of `Database`, `SyncAsyncQuery` and `DatabaseInteraction` interfaces.<br>
Object for accessing the database.<br>
It can send both synchronous and asynchronous queries.

__IMPORTANT__: after the end of work it is necessary to close the connection using Close method.

For each transaction a new instance of the `Transaction` object is created,
so each transaction is executed in its own scope and is completely safe.
```golang
type MysqlDatabase struct {
	username string
	password string
	host     string
	port     string
	database string
	db       *sql.DB
	tx       *sql.Tx
	syncQ    interfaces.SyncQ
	asyncQ   interfaces.AsyncQ
}
```

#### MysqlDatabase.Open
Connecting to a mysql database.<br>
Also, initialization of synchronous and asynchronous queries.
```golang
func (d *MysqlDatabase) Open() error {
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
```

#### MysqlDatabase.Close
Closes the connection to the database.
```golang
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
```

#### MysqlDatabase.NewTransaction
```golang
func (d *MysqlDatabase) NewTransaction() (interfaces.DatabaseTransaction, error) {
	return NewMysqlTransaction(d.db, d.syncQ, d.asyncQ)
}
```

#### MysqlDatabase.SyncQ
Getting access to synchronous requests.
```golang
func (d *MysqlDatabase) SyncQ() interfaces.SyncQ {
	return d.syncQ
}
```

#### MysqlDatabase.NewAsyncQ
Creates and returns a new instance of `interfaces.AsyncQ`.<br>
This is necessary for data security, since `interfaces.AsyncQ` stores SQL query data, so a separate instance must be 
created for each HTTP handler.
```golang
func (d *MysqlDatabase) NewAsyncQ() (interfaces.AsyncQ, error) {
	aq, err := d.asyncQ.New()
	if err != nil {
		return nil, err
	}
	return aq.(interfaces.AsyncQ), nil
}
```

### MysqlTransaction
An object that performs transactions to the mysql database.<br>
This object is used only for one transaction, for each next transaction a new instance of the object must be created.
```golang
type MysqlTransaction struct {
	db     *sql.DB
	tx     *sql.Tx
	syncQ  interfaces.SyncQ
	asyncQ interfaces.AsyncQ
}
```

#### MysqlTransaction.BeginTransaction
Starts the transaction.<br>
Only one transaction can be started per object instance.
```golang
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
```

#### MysqlTransaction.CommitTransaction
Writes the transaction to the database.
```golang
func (t *MysqlTransaction) CommitTransaction() error {
	if t.tx == nil {
		return errors.New("transaction not begin")
	}
	if err := t.tx.Commit(); err != nil {
		return err
	}
	t.tx = nil
	return nil
}
```

#### MysqlTransaction.RollBackTransaction
Undoes any changes that were made during the transaction.<br>
That is, after executing the `BeginTransaction` method.
```golang
func (t *MysqlTransaction) RollBackTransaction() error {
	if t.tx == nil {
		return errors.New("transaction not begin")
	}
	if err := t.tx.Rollback(); err != nil {
		return err
	}
	t.tx = nil
	return nil
}
```

#### MysqlTransaction.SyncQ
Getting access to synchronous requests.
```golang
func (t *MysqlTransaction) SyncQ() interfaces.SyncQ {
	return t.syncQ
}
```

#### MysqlTransaction.NewAsyncQ
Creates and returns a new instance of `interfaces.AsyncQ`.<br>
This is necessary for data security, since `interfaces.AsyncQ` stores SQL query data, so a separate instance must be 
created for each HTTP handler.
```golang
func (t *MysqlTransaction) NewAsyncQ() (interfaces.AsyncQ, error) {
	aq, err := t.asyncQ.New()
	if err != nil {
		return nil, err
	}
	return aq.(interfaces.AsyncQ), nil
}
```

### DbQuery
Standard database queries. They are used *sql.DB.
Requests are executed as usual.
```golang
type DbQuery struct {
	DB *sql.DB
}
```

#### DbQuery.Query
Used to execute queries that return data.
For example, the SELECT command.
```golang
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
```

#### DbQuery.Exec
Used to execute queries that do not return data.
For example, the INSERT command.

Returns the following data:

* Key "insertID" is the identifier of the inserted row using INSERT.
* Key "rowsAffected" - Returns the number of rows affected by INSERT, UPDATE, DELETE.
```golang
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
```

### DbTxQuery
Queries that can be rolled back. Used *sql.Tx.<br>
This object will perform queries with the `*sql.Tx` object that is used for transactions.
```golang
type DbTxQuery struct {
	Tx *sql.Tx
}
```

#### DbTxQuery.Query
Used to execute queries that return data.
For example, the SELECT command.
```golang
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
```

#### DbTxQuery.Exec
Used to execute queries that do not return data.
For example, the INSERT command.

Returns the following data:

* Key "insertID" is the identifier of the inserted row using INSERT.
* Key "rowsAffected" - Returns the number of rows affected by INSERT, UPDATE, DELETE.
```golang
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
```

___

#### InitDatabasePool
Initializes a database pool.<br>
Only one pool is created, which is specified in the `Default.Database.MainConnectionPoolName` settings.
Once created, the pool is locked. Therefore, you need to initialize it manually if you need more connections.
```golang
func InitDatabasePool(manager interfaces.Manager, db interfaces.Database) error {
	name := config.LoadedConfig().Default.Database.MainConnectionPoolName
	if err := manager.Database().AddConnection(name, db); err != nil {
		return err
	}
	manager.Database().Lock()
	return nil
}
```