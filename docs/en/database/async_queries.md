## IAsyncQueries
The difference between this interface and ISyncQueries is that here requests are sent asynchronously. Therefore, they will not be listed here 
query methods, because they work identically only asynchronously. It's also worth noting that this interface depends on ISyncQueries, 
because it directly uses its query methods.<br>
Each query method has a parameter ``key string``, this parameter sets the key for the result of the query which is then executed 
can be used in the __LoadAsyncRes__ method.

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

__SetSyncQueries__
```
SetSyncQueries(queries interfaces.ISyncQueries)
```
Sets the ISyncQueries interface to access for its query methods.

__Wait__
```
Wait()
```
The method waits for all asynchronous requests to complete. It should always be used.

__LoadAsyncRes__
```
LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
```
The method waits for all asynchronous requests is complete. It should always be usd.
__IMPORTANT:__ this method should only be used after the __Wait__ method.
