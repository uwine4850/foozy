## ISyncQueries
Перед описом методів варто зазначити, що кожен метод який робить запит повертає ``[]map[string]interface{}`` - ключ дорівнює 
назві стовпця, а __interface{}__ дорівнює його значенню(якщо значень немає, буде порожня карта). Ці значення можна конвертувати самостійно, або функціями із 
пакету [dbutils](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/dbutils/dbutils.md).

Детальніше про взаємодію із базою даних можна подивитись у цих [тестах](https://github.com/uwine4850/foozy/tree/master/tests/dbtest).

__Query__
```
Query(query string, args ...any) ([]map[string]interface{}, error)
```
Відправляє параметризований запит до бази даних.

__SetDB__
```
SetDB(db *sql.DB)
```
Встановлює sql підключення до бази даних.

__Select__
```
Select(rows []string, tableName string, where []dbutils.DbEquals, limit int) ([]map[string]interface{}, error)
```
Виконує sql запит SELECT. Параметр ``rows`` це стовпці які будуть виводитися(* - всі). Параметр ``where`` це масив структур
dbutils.DbEquals де ключ це стовпець, а значення це значення цього стовпця.

__Insert__
```
Insert(tableName string, params map[string]interface{}) ([]map[string]interface{}, error)
```
Виконує sql запит INSERT. Параметр ``params`` це дані для вставки, а саме ключ дорівнює стовпцю, а інтерфейс дорівнює його значенню.

__Delete__
```
Delete(tableName string, where []dbutils.DbEquals) ([]map[string]interface{}, error)
```
Виконує sql запит DELETE. Параметр ``where`` це масив структур dbutils.DbEquals де ключ це стовпець, а значення це 
значення цього стовпця. Тобто, метод видаляє всі стовпці які підходять під умову where.

__Update__
```
Update(tableName string, params []dbutils.DbEquals, where []dbutils.DbEquals) ([]map[string]interface{}, error)
```
Виконує sql запит UPDATE. Параметр ``params`` це масив структур dbutils.DbEquals де ключ це стовпець, а значення це нове
значення цього стовпця. Параметр ``where`` це масив структур dbutils.DbEquals який відповідає за умову.