## package qb
This package is designed for creating database queries.<br>
Basic use of the package:
* Open a database connection using the __database__ package.
* Create a new instance of a `QB` object.
* Create a chain of sql commands.
* When the chain of commands is complete, merge it using the `Merge` method.
* After all these operations, execute the `Query` or `Exec` method. These methods send a query to the database.


### type QB struct
An object for creating querybuild queries.<br>
To execute a query, perform the following steps:
* Create a chain of commands using the corresponding methods and merge them using the `Merge` method.
* Using `Query` or `Exec` method to execute sql command.

The sql command corresponds to the name of the method.

### type customFunc struct
The `customFunc` object adds the ability to use `subquery` anywhere in the query. This allows you to make more flexible queries.<br>
To use it, you just need to call the `Func` method anywhere and pass a `subquery` instance to it.

### type subquery struct
sql subquery that is inserted in the middle of the main query.
The bracket field determines whether this query will be bracketed.
It is important to specify that a new instance of the `QB` object is used.
The `QB` object can also be without an initialized database, since it is not used.
The `QB` structure is used only to receive the sql string and arguments, the query itself is not sent.

__ParseSubQuery__
```
ParseSubQuery(sq any, outQueryString *string, outQueryArgs *[]any) bool
```
Processes an instance of subQuery.<br>
outQueryString — sql string of the subquery.<br>
outQueryArgs — position arguments of sql subquery.

## Filters
Here we will describe filters for commands such as __WHERE__ and the like. The name of methods corresponds to the name of commands.

### type compare struct
Subject to comparative conditions. Used in some methods:
* Where
* Having
* InnerJoin
* LeftJoin
* RightJoin

This object takes three values: left operand, operator, and right operand.
Example Usage:
```go
qb.Compare("id", qb.EQUAL, 1)
```

### type noArgsCompare struct
noArgsCompare does almost everything that a default `compare` object does.
The peculiarity is that this object does not pass the right operand as a
positional argument. So the "?" sign will not be used.
It is important to clarify that the structure still processes and
passes arguments of external objects like `subQuery`.

### type between struct
Builds a sql filter BETWEEN.<br>
Uses three arguments: column name, left operand and right operand.<br>
Example usage:
```go
qb.Between("id", 1, 10)
```

### type array struct
Makes a sql array. For example, (1, 3, 5).<br>
Example usage:
```go
qb.Compare("id", qb.IN, qb.Array(1, 3, 4))
```

### type exists struct
Builds a sql filter EXISTS.<br>
Example usage:
```go
...Select(Exists(exSQ))
```

## Functions for querybuild.
__SelectExists__
```
SelectExists(qb *QB, tableName string, whereValues ...any) (bool, error)
```
SelectExists checks if there is a value in the table. It is important to use a condition for correct operation.<br>
Example of use:
```go
ex, err := qb.SelectExists(qb.NewSyncQB(db.SyncQ()), "roles",
		qb.Compare("user_id", qb.EQUAL, userID), qb.AND,
		qb.Compare("role_name", qb.EQUAL, roleName)
)
```

__Increment__
```
Increment(qb *QB, tableName string, field string, whereValues ...any) (bool, error)
```
Increases the numeric value of the table by one.<br>
Example of use:
```go
qb.Increment(qb.NewSyncQB(db.SyncQ()), "table", "id", qb.Compare("value", qb.EQUAL, 1))
```