## Package dbutils
Helper functions and structures for the [database package](https://github.com/uwine4850/foozy/blob/master/docs/en/database/database.md).

You can see more about interaction with the database in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/dbtest).

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
Structure of general use. The ``name`` field is the name of the column, and the ``value`` field is its value.

__RepeatValues__
```
RepeatValues(count int, sep string) string
```
It is used for parameterized queries, namely for repeating the ``?`` sign.

__ScanRows__
```
ScanRows(rows *sql.Rows, fn func(row map[string]interface{}))
```
Reads the result of ``*sql.Rows``, this data type which contains the values ​​of several (both one) rows.<br>
The task of this method is to read each line and convert it to the format __map[string]interface{}__ where the key is the name of the column, 
and interface{} is its value. And the last task of this function is to run the ``fn`` method for each iteration.

__ParseParams__
```
ParseParams(params map[string]interface{}) ([]string, []interface{})
```
Converts a map into two results of type []string and []interface{}, where the first is the keys and the second is the value of the keys.

__ParseEquals__
```
ParseEquals(equals []DbEquals, conjunction string) (string, []interface{})
```
Converts ``equals []DbEquals`` into a string for a parameterized query, namely for sql code where there is a ``=`` sign.
It is also possible to set the required separator<br>.
For example, the string can be ``key1 = ?, key2 = ?``. In addition, the ``[]interface{}`` value, which contains an array, is returned 
key values.

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
Converts the value from the query result to a date and time.

__ParseFloat__
```
ParseFloat(value interface{}) (float64, error)
```
Converts the value from the query result to a comma-separated number.

__FillStructFromDb__
```
FillStructFromDb(dbRes map[string]interface{}, fill interface{}) error
```
Fills the structure with data from the database.
Each variable of the placeholder structure must have a "db" tag, which is responsible for the name of the column in
database, for example `db: "name"`.

__FillMapFromDb__
```
FillMapFromDb(dbRes map[string]interface{}, fill *map[string]string) error
```
Fills the map with data from the database.

__FillReflectValueFromDb__
```
FillReflectValueFromDb(dbRes map[string]interface{}, fill *reflect.Value) error
```
Fills a structure whose type is *reflect.Value. That is, the method fills the data from the database into the structure, which is created with the help package reflect.