## Package auth
У цьому пакеті реалізована базова логіка аутентифікації. Більш детальний алгоритм аутентифікації фожна побачити у [цьому тесті](https://github.com/uwine4850/foozy/blob/master/tests/authtest/auth_test.go).

Алгоритм аутентифікації полягає у наступному:
* Створення таблиці аутентифікації з допомогою метода __CreateAuthTable__.
* Реєстрація користувача з допомогою метода __RegisterUser__.
* Вхід користувача методом __LoginUser__.
* Для того, щоб ключі кодування оновлювались потрібно використовувати __globalflow__ та функцію *bglobalflow.KeyUpdater(1)*.
    ```
    gf := globalflow.NewGlobalFlow(1)
    gf.AddNotWaitTask(bglobalflow.KeyUpdater(1))
    gf.Run(manager)
    ```
* Тепер, коли ключі оновлюються, потрібно оновити кодування з допомогою Middlewares. Це можна зробити з допомогою метода *builtin_mddl.Auth*.
    ```
    mddl := middlewares.NewMiddleware()
    mddl.HandlerMddl(0, builtin_mddl.Auth("/login", mddlDb))
    ```

Важливо зазначити, що пропускаючи ітерацію globalflow не оновивши кодування, користувач не зможе їх оновити і потрібно буде перезаходити в аккаунт. Це відбувається тому, що збегігаються тільки передостанні ключі, а для зміни кодування потрібно мати підходящі ключі. Відповідно, якщо пропустити ітерацію, потрібні ключі для кодування просто будуть втрачені.

## Методи структури Auth

__type AuthCookie struct__
Структура призначена для опису даних cookie.

__type User struct__
Структура призначена для опису даних користувача.

__RegisterUser__
```
RegisterUser(username string, password string) error
```
Створює нового користувача. Якщо щось пішло не так повертає 
помилку.

__LoginUser__
```
LoginUser(username string, password string) error
```
Перевіряє правильність даних користувача. Якщо немає помилки - всі дані коректні. 
Повертає структуру __User__, яка містить в собі дані про користувача. Також 
додає дані cookie про користувача.

__UpdateAuthCookie__
```
UpdateAuthCookie(hashKey []byte, blockKey []byte, r *http.Request) error
```
Оновлює шифрування даних користувача, які знаходяться у cookie. Для вдалого оновлення потрібно у поля *hashKey* та *blockKey* передавати ключі які використовувалися для шифрування.

__ChangePassword__
```
ChangePassword(username string, oldPassword string, newPassword string) error
```
Для зміни пароля потрібно ввести ім'я користувача, старий пароль та новий пароль.

__UserByUsername__
```
UserByUsername(username string) (map[string]interface{}, error)
```
Перевіряє чи є користувач в базі даних. Якщо він є, повертає його дані.

__UserById__
```
UserById(id any) (map[string]interface{}, error)
```
Шукає користувача по ID. Якщо він є, повертає його дані.

## Глобальні функції для пакета auth

__CreateAuthTable__
```
CreateAuthTable(database interfaces.IDatabase) error
```
Створює таблицю користувачів в базі даних. Використовується в конструкторі структури ``Auth``.

__HashPassword__
```
HashPassword(password string) (string, error)
```
Хешує пароль для зберігання.

__ComparePassword__
```
ComparePassword(hashedPassword string, password string) error
```
Порівнює чи збігається хеш пароль і пароль.

## Приклад використання
```
dbArgs := database.DbArgs{
	Username: "root", Password: "1111", Host: "localhost", Port: "3408", DatabaseName: "foozy_test",
}
db := database.NewDatabase(dbArgs)
err := db.Connect()
if err != nil {
    panic(err)
}
defer func(db *database.Database) {
    err := db.Close()
    if err != nil {
	    panic(err)
	}
}(db)
_auth := auth.NewAuth(db)
err = _auth.RegisterUser("user", "111111")
if err != nil {
    fmt.Println(err)
}
_, err := _auth.LoginUser("user", "111111")
if err != nil {
    fmt.Println(err)
	return
}
```