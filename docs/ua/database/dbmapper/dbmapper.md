## package dbmapper
Записує дані у вибраний об'єкт. Може проводити різноманітні операції над 
заповненими даними.

Подивитись на роботу пакету можна у цих [тестах](https://github.com/uwine4850/foozy/tree/master/tests/dbtest/dbmapper_test).

### type Mapper struct
Структура, яка використовується для того, щоб заповнити об'єкт данними із бази даних.
У `DatabaseResult` передається список із данних із БД. У поле `Output` потрібно 
передати посилання на потрібний об'єкт. 

Важливо зазничити, що дані повині бути у вигляді зрізу.

На даний момент у полі `Output` можливо використовувати три вида даних:
* Структура.
* Структура у викляді `reflect.Value`.
* Тип даних map[string]string.

Якщо як `Output` використовується структура, її потрібно правильно налаштувати.
Приклад структури:
```
type DbTestMapper struct {
	Col1 string `db:"col1"`
	Col2 string `db:"col2"`
	Col3 string `db:"col3" empty:"-err"`
	Col4 string `db:"col4" empty:"0"`
}
```
Для налаштування можливо використовувати два тега:
* db — обов'язковий тег. Потрібний для того щоб знати яке поле відповідає якому 
стовпчику у таблиці. Відповідно, у ньому потрібно вказати назву стовпця.
* empty — не обов'язковий тег. Цей тег використовується для команди того, що 
потрібно робити, коли поле порожнє. На даний момент є тільки два значення:
    * -err — якщо поле порожнє, виводить відповідну помилку.
    * просто текст — текст який замінить пусті значення.

__Fill__
```
Fill() error 
```
Заповнює данні.

### Інші функції пакету.

__FillStructFromDb__
```
FillStructFromDb(dbRes map[string]interface{}, fillPtr itypeopr.IPtr) error
```
Заповнює структуру даними з бази даних.
Кожна змінна заповнюваної структури повинна мати тег "db", який відповідає за назву 
стовпця в базі даних, наприклад, `db: "name"`.
Також можна використати тег `empty` про який детальніше написана вище.

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
ParamsValueFromStruct(filledStructurePtr itypeopr.IPtr, nilIfEmpty []string) (map[string]any, error)
```
Створює карту зі структури, яка описує таблицю.
Для правильної роботи вам потрібна завершена структура, а обов’язкові поля мають мати тег `db:"<назва стовпця>"`.
Також можна використати тег `empty` про який детальніше написана вище.