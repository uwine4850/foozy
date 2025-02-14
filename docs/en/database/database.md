## Database package
This package is required for convenient use of the database. This package does not depend on others and they do not depend on it, therefore 
it is possible to use this if necessary.<br>
The database consists of several important interfaces:
* IDatabase - connection and disconnection from the database. You can also use this interface to access 
request interfaces.
* ISyncQueries - an interface for sending synchronous requests to the database.
* IAsyncQueries - an interface for sending asynchronous requests to the database.

You can see more about interaction with the database in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/database/db_test).

## IDatabase
Database structure for accessing the database.
It can send both synchronous and asynchronous queries.
IMPORTANT: after the end of work it is necessary to close the connection using Close method.

Principle of `db` and `tx` swapping:<br>
After the initial initialization, `db` is used. These are standard database queries. If the `BeginTransaction` method is used,
the `db` instance will be replaced by the `tx` instance. Then the `CommitTransaction` method changes them back.
These instances use the same interface, so they are directly used by the `ISyncQueries` interface, which in turn is used by
the `IAsyncQueries` interface.
The main difference between `db` and `tx` objects is that the latter is used for the ability to cancel database transactions.
__Connect__
```
Connect() error
```
Database connection. Initializes the ISyncQueries and IAsyncQueries interfaces.<br>
__IMPORTANT:__ After you finish working with the database, you need to disconnect from it using the ``Close`` method.

__Ping__
```
Ping() error
```
Checking the connection to the database.

__BeginTransaction__
```
BeginTransaction()
```
Transaction execution begins.
Changes the executable object for database queries, so all queries executed by this method will use *sql.Tx.

__CommitTransaction__
```
CommitTransaction() failed.
```
Writes changes to the database.
This method ends the transaction and changes the execution of the object query from *sql.Tx to *sql.DB.
Therefore, all subsequent queries are applied using *sql.DB.

__RollBackTransaction__
```
RollBackTransaction() error
```
Returns changes made by commands AFTER `BeginTransaction()` is called.
This method ends the transaction and changes the executing query object from *sql.Tx to *sql.DB.

__Close__
```
Close() error
```
Disconnecting from the database.

__SetSyncQueries__
```
SetSyncQueries(q interfaces.ISyncQueries)
```
Establishes the synchronous query interface for accessing them from IDatabase.

__SetAsyncQueries__
```
SetAsyncQueries(q interfaces.IAsyncQueries)
```
Establishes an asynchronous query interface for accessing them from IDatabase.

__SyncQ__
```
SyncQ() interfaces.ISyncQueries
```
Access to synchronous requests.

__AsyncQ__
```
AsyncQ() interfaces.IAsyncQueries
```
Access to asynchronous requests.

__DatabaseName__
```
DatabaseName() string
```
Returns the name of the database.