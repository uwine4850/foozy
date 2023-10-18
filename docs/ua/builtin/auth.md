## Package auth
У цьому пакеті реалізована базова логіка аутентифікації.

## Методи структури Auth
__RegisterUser__
```
RegisterUser(username string, password string) error
```
Створює нового користувача. Якщо щось пішло не так повертає помилку.

__LoginUser__
```
LoginUser(username string, password string) error
```
Перевіряє правильність даних користувача. Якщо немає помилки - всі дані коректні.

__ChangePassword__
```
ChangePassword(username string, oldPassword string, newPassword string) error
```
Для зміни пароля потрібно ввести ім'я користувача, старий пароль та новий пароль.

__UserExist__
```
UserExist(username string) (map[string]interface{}, error)
```
Перевіряє чи є користувач в базі даних. Якщо він є, повертає його дані.

## Глобальні пакети для пакета auth

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
db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
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
_auth, err := auth.NewAuth(db)
if err != nil {
    panic(err)
}
err = _auth.RegisterUser("user", "111111")
if err != nil {
    fmt.Println(err)
}
err = _auth.LoginUser("user", "111111")
if err != nil {
    fmt.Println(err)
	return
}
```