## ISyncQueries
Before describing the methods, it should be noted that each method that makes a request returns ``[]map[string]interface{}`` - the key is 
name of the column, and __interface{}__ is equal to its value (if there are no values, there will be an empty map). These values ​​can be converted independently or with functions from the [dbutils package](https://github.com/uwine4850/foozy/blob/master/docs/en/database/dbutils/dbutils.md).

You can see more about interaction with the database in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/database/db_test).

__Query__
```
Query(query string, args ...any) ([]map[string]interface{}, error)
```
Sends a parameterized query to the database.

__SetDB__
```
SetDB(db *sql.DB)
```
Establishes a sql connection to the database.

__Select__
```
Select(rows []string, tableName string, where []dbutils.DbEquals, limit int) ([]map[string]interface{}, error)
```
Executes sql SELECT query. The ``rows`` parameter is the columns that will be displayed (* - all). The ``where`` parameter is an array of structures dbutils.DbEquals where key is a column and value is the value of that column.

__Insert__
```
Insert(tableName string, params map[string]any) ([]map[string]interface{}, error)
```
Executes sql INSERT query. The ``params`` parameter is the data to insert, namely the key is equal to the column and the interface is equal to its value.

__Delete__
```
Delete(tableName string, where []dbutils.DbEquals) ([]map[string]interface{}, error)
```
Executes a sql DELETE query. The ``where`` parameter is an array of dbutils.DbEquals structures where the key is the column and the value is this 
the value of this column. That is, the method removes all columns that fit the where condition.

__Update__
```
Update(tableName string, params map[string]any, where []dbutils.DbEquals) ([]map[string]interface{}, error)
```
Executes an UPDATE sql query. The ``params`` parameter is the data to update, where the key is equal to the column name and the map value 
is the new value of the selected column. The ``where`` parameter is an array of dbutils.DbEquals structures that is responsible for the 
condition.