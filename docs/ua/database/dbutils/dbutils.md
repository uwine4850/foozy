## Package dbutils
Допоміжні функції та структури для пакета [database](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/database.md).

Детальніше про взаємодію із базою даних можна подивитись у цих [тестах](https://github.com/uwine4850/foozy/tree/master/tests/dbtest).

__AsyncQueryData__
```
type AsyncQueryData struct {
    Res   []map[string]interface{}
    Error string
}
```
Структура, яка використовується для виводу результатів асинхронного sql запиту.

__DbEquals__
```
type DbEquals struct {
    Name  string
    Value interface{}
}
```
Структура загального використання. Поле ``name`` це ім'я стовпця, а поле ``value`` це його значення.

__RepeatValues__
```
RepeatValues(count int, sep string) string
```
Використовується для параметризованих запитів, а саме для повтору знака ``?``.

__ScanRows__
```
ScanRows(rows *sql.Rows, fn func(row map[string]interface{}))
```
Читає результат ``*sql.Rows``, цей тип даних який містить значення декількох(обо одного) рядків.<br>
Задання цього методу - зчитати кожен рядок та конвертувати його в формат __map[string]interface{}__ де ключ це назва колонки, 
а interface{} це його значення. І останнє завдання цієї функції - це запустити метод ``fn`` для кожної ітерації.

__ParseParams__
```
ParseParams(params map[string]interface{}) ([]string, []interface{})
```
Перетворює карту на два результати типу []string та []interface{}, де перший це ключі, а другий це значення ключів.

__ParseEquals__
```
ParseEquals(equals []DbEquals, conjunction string) (string, []interface{})
```
Перетворює ``equals []DbEquals`` в рядок для параметризованого запиту, а саме для sql коду де є знак ``=``. 
Також є можливість встановити потрібний роздільник<br>.
Наприклад, рядок може бути такий ``key1 = ?, key2 = ?``. Крім цього повертається значення ``[]interface{}`` яке містить масив 
значень ключів.

__ParseString__
```
ParseString
```
Перетворює значення із результату запиту в рядок.

__ParseInt__
```
ParseInt(value interface{}) (int, error)
```
Перетворює значення із результату запиту в ціле число.

__ParseDateTime__
```
ParseDateTime(layout string, value interface{}) (time.Time, error)
```
Перетворює значення із результату запиту в дату та час.

__ParseFloat__
```
ParseFloat(value interface{}) (float64, error)
```
Перетворює значення із результату запиту в число з комою.

__FillStructFromDb__
```
FillStructFromDb(dbRes map[string]interface{}, fill interface{}) error
```
Заповнює структуру даними з бази даних.
Кожна змінна заповнюваної структури повинна мати тег "db", який відповідає за назву стовпця в
базі даних, наприклад, `db: "name"`.

__FillMapFromDb__
```
FillMapFromDb(dbRes map[string]interface{}, fill *map[string]string) error
```
Заповнює мапу даними із бази даних.

__FillReflectValueFromDb__
```
FillReflectValueFromDb(dbRes map[string]interface{}, fill *reflect.Value) error
```
Заповнює структуру тип якої *reflect.Value. Тобто, метод заповнює дані із бази даних у структуру, яка стрворена з допомогою 
пакету reflect.

__ParamsValueFromStruct__
```
ParamsValueFromStruct(structure interface{}) (map[string]any, error)
```
Створює карту зі структури, яка описує таблицю.
Для правильної роботи вам потрібна завершена структура, а обов’язкові поля мають мати тег `db:"<назва стовпця>"`.