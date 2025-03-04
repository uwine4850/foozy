## ISyncQueries
Before describing the methods, it should be noted that each method that makes a request returns ``[]map[string]interface{}`` - the key is 
name of the column, and __interface{}__ is equal to its value (if there are no values, there will be an empty map). These values ​​can be converted independently or with functions from the [dbutils package](https://github.com/uwine4850/foozy/blob/master/docs/en/database/dbutils/dbutils.md).

You can see more about interaction with the database in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/database/db_test).

__Query__
```
Query(query string, args ...any) ([]map[string]interface{}, error)
```
Sends a parameterized query to the database.

__Exec__
```
Exec(query string, args ...any) (map[string]interface{}, error)
```
Used to execute queries that do not return a result, such as UPDATE or INSERT.
The method returns a map containing two values:
* insertID — returns the ID of the field that was inserted using the INSERT command. It is important that there is an AUTO INCREMENT field.
* rowsAffected — the number of rows that were affected during the query, for example, during the UPDATE command.

__SetDB__
```
SetDB(db *sql.DB)
```
Establishes a sql connection to the database.
