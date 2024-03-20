## Database package
This package is required for convenient use of the database. This package does not depend on other packages, and they do not depend on it, so
it is possible to use this package if necessary.<br>
The database consists of several important interfaces:
* IDatabase - connecting and disconnecting from the database. You can also use this interface to access
  query interfaces.
* ISyncQueries - interface for sending synchronous queries to the database.
* IAsyncQueries - interface for sending asynchronous queries to the database.

These interfaces will be described below.

## IDatabase
__Connect__
```
Connect() error
```
Connects to the database. Initializes the ISyncQueries and IAsyncQueries interfaces.<br>
__IMPORTANT:__ After you finish working with the database, you need to disconnect from it using the ``Close`` method.

```
Ping() error
```
Checking the connection to the database.

__Close__
```
Close() error
```
Disconnecting from the database.

__SetSyncQueries__
```
SetSyncQueries(q interfaces.ISyncQueries)
```
Sets the interface of synchronous queries to access them from IDatabase.

__SetAsyncQueries__
```
SetAsyncQueries(q interfaces.IAsyncQueries)
```
Sets the asynchronous queries interface for accessing them from the IDatabase.

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

## ISyncQueries
Before describing the methods, it should be noted that each method that makes a request returns ``[]map[string]interface{}`` - the key is equal to
the column name, and __interface{}__ is equal to its value (if there are no values, there will be an empty map). These values can be converted independently or with functions from the
the [dbutils package](https://github.com/uwine4850/foozy/blob/master/docs/en/dbutils.md).

__Query__
```
Query(query string, args ...any) ([]map[string]interface{}, error)
```
Sends a parameterized query to the database.

__SetDB__
```
SetDB(db *sql.DB)
```
Встановлює sql підключення до бази даних.

__Select__
```
Select(rows []string, tableName string, where []dbutils.DbEquals, limit int) ([]map[string]interface{}, error)
```
Executes the sql query SELECT. The ``rows`` parameter is the columns that will be displayed (* - all). The 'where' parameter is an array of structures
dbutils.DbEquals where the key is a column, and the value is the value of this column.

__Insert__
```
Insert(tableName string, params map[string]interface{}) ([]map[string]interface{}, error)
```
Executes the sql query INSERT. The ``params`` parameter is the data to be inserted, namely the key is equal to the column, 
and the interface is equal to its value.
__Delete__
```
Delete(tableName string, where []dbutils.DbEquals) ([]map[string]interface{}, error)
```
Executes the sql query DELETE. The ``where`` parameter is an array of dbutils.DbEquals structures where the key is a column, and the value is
is the value of this column. That is, the method deletes all columns that match the where condition.

__Update__
```
Update(tableName string, params []dbutils.DbEquals, where []dbutils.DbEquals) ([]map[string]interface{}, error)
```
Executes the sql query UPDATE. The ``params`` parameter is an array of dbutils.DbEquals structures where the key is a column, and the value is a new
the value of this column. The ``where`` parameter is an array of dbutils.DbEquals structures that is responsible for the condition.

## IAsyncQueries
The difference between this interface and ISyncQueries is that requests are sent asynchronously. Therefore, we will not list the
query methods, because they work identically only asynchronously. It is also worth noting that this interface depends on ISyncQueries,
because it directly uses its query methods.<br>
Each query method has a ``key string`` parameter, this parameter sets the key for the result of the query, which can then
can be used in the __LoadAsyncRes__ method.

__SetSyncQueries__
```
SetSyncQueries(queries interfaces.ISyncQueries)
```
Sets the ISyncQueries interface to be accessed by its query methods.

__Wait__
```
Wait()
```
The method waits for all asynchronous requests to complete. It should always be used.

__LoadAsyncRes__
```
LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
```
Print the result of executing the query by the key that was set earlier.<br>
__IMPORTANT:__ This method should only be used after the __Wait__ method.
