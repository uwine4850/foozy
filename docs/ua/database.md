## Database package
Даний пакет потрібен для зручного користування базою даних. Цей пакет не залежить від інших як і вони від нього, тому 
є можливість використовувати даний за необхідності.<br>
База даних складається з декількох важливих інтерфейсів:
* IDatabase - підключення та відключення від бази даних. Також з допомогою цього інтерфейсу можна отримати доступ до 
інтерфейсів запитів.
* ISyncQueries - інтерфейс для відправки синхронних запитів до бази даних.
* IAsyncQueries - інтерфейс для відправки асинхронних запитів до бази даних.

Далі будуть описуватися ці інтерфейси.

## IDatabase
__Connect__
```
Connect() error
```
Підключення до бази даних. Ініціалізує інтерфейси ISyncQueries та IAsyncQueries.<br>
__ВАЖЛИВО:__ Після завершення роботи з базою даних потрібно відключитись від неї за допомогою метода ``Close``.

__Close__
```
Close() error
```
Відключення від бази даних.

__SetSyncQueries__
```
SetSyncQueries(q interfaces.ISyncQueries)
```
Встановлює інтерфейс синхронних запитів для доступу до них з IDatabase.

__SetAsyncQueries__
```
SetAsyncQueries(q interfaces.IAsyncQueries)
```
Встановлює інтерфейс асинхронних запитів для доступу до них з IDatabase.

__SyncQ__
```
SyncQ() interfaces.ISyncQueries
```
Доступ до синхронних запитів.

__AsyncQ__
```
AsyncQ() interfaces.IAsyncQueries
```
Доступ до асинхронних запитів.

__DatabaseName__
```
DatabaseName() string
```
Повертає назву бази даних.

## ISyncQueries
Перед описом методів варто зазначити, що кожен метод який робить запит повертає ``[]map[string]interface{}`` - ключ дорівнює 
назві стовпця, а __interface{}__ дорівнює його значенню(якщо значень немає, буде порожня карта). Ці значення можна конвертувати самостійно, або функціями із 
пакету [dbutils](https://github.com/uwine4850/foozy/blob/master/docs/ua/dbutils.md).

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

## IAsyncQueries
Відмінність цього інтерфейсу від ISyncQueries це те, що тут запити відправляються асинхронно. Тому тут не будуть перечислені 
методи запитів, адже вони працюють ідентично тільки асинхронно. Також варто зазначити, що цей інтерфейс залежить від ISyncQueries, 
адже він напряму використовує його методи запитів.<br>
Кожен метод запитів має параметр ``key string``, цей параметр встановлює ключ для результату виконання запита який потім 
можна використати у методі __LoadAsyncRes__.

__SetSyncQueries__
```
SetSyncQueries(queries interfaces.ISyncQueries)
```
Встановлює інтерфейс ISyncQueries для доступу для його методів запиту.

__Wait__
```
Wait()
```
Метод чекає завершення усіх асинхронних запитів. Його потрібно використовувати завжди.

__LoadAsyncRes__
```
LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool)
```
Вивід результату виконання запиту за ключем який був встановлений раніше.<br>
__ВАЖЛИВО:__ цей метод потрібно використовувати лише після методу __Wait__.
