## Package dbutils
Auxiliary functions and structures for the [database package](https://github.com/uwine4850/foozy/blob/master/docs/en/database.md).

__AsyncQueryData__
```
type AsyncQueryData struct {
    Res   []map[string]interface{}
    Error string
}
```
A structure used to display the results of an asynchronous sql query.

__DbEquals__
```
type DbEquals struct {
    Name  string
    Value interface{}
}
```
The structure of general use. The ``name`` field is the name of the column, and the ``value`` field is its value.
__RepeatValues__
```
RepeatValues(count int, sep string) string
```
It is used for parameterized queries, namely for repeating the ``?`` sign.

__ScanRows__
```
ScanRows(rows *sql.Rows, fn func(row map[string]interface{}))
```
Reads the result of ``*sql.Rows``, this data type that contains the values of several (or one) rows.<br
The task of this method is to read each row and convert it to the format __map[string]interface{}__ where the key is the name of the column,
and interface{} is its value. And the last task of this function is to run the ``fn`` method for each iteration.

__ParseParams__
```
ParseParams(params map[string]interface{}) ([]string, []interface{})
```
Converts the map to two results of type []string and []interface{}, where the first is the keys, and the second is the key values.
__ParseEquals__
```
ParseEquals(equals []DbEquals, conjunction string) (string, []interface{})
```
Converts ``equals []DbEquals`` to a string for a parameterized query, namely for sql code where there is a ``=`` sign.
It is also possible to set the desired separator<br>.
For example, the string could be ``key1 = ?, key2 = ?``. In addition, the value ``[]interface{}`` is returned, which contains an array
of key values.

__ParseString__
```
ParseString
```
Converts the value from the query result to a string.

__ParseInt__
```
ParseInt(value interface{}) (int, error)
```
Converts the value from the query result to an integer.

__ParseDateTime__
```
ParseDateTime(layout string, value interface{}) (time.Time, error)
```
Converts values from the query result to date and time.

__ParseFloat__
```
ParseFloat(value interface{}) (float64, error)
```
Converts the value from the query result to a decimal number.

__FillStructFromDb__
```
FillStructFromDb(dbRes map[string]interface{}, fill interface{}) error
```
Fills the structure with data from the database.
Each variable of the filled structure must have a "db" tag, which is responsible for the name of the column in the
the database, for example, `db: "name"`.
