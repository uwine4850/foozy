## Form package
Пакет form містить весь доступний функціонал(на даний момент) для роботи з формами.

## Методи які використовує інтерфейс IForm
__Parse__
```
Parse() error
```
Метод який повинен завжди використовуватися для парсингу форми. З його допомогою виконується базова обробка форми.

__GetMultipartForm__
```
GetMultipartForm() *multipart.Form
```
Повертає дані влаштовану структуру golang ``multipart.Form``. Цей метод використовується для форми ``multipart/form-data``.

__GetApplicationForm__
```
GetApplicationForm() url.Values
```
Повертає дані влаштовану структуру golang ``url.Values``. Цей метод використовується для форми ``application/x-www-form-urlencoded``.

__Value__
```
Value(key string) string
```
Повертає дані из текстового поля форми.

__File__
```
File(key string) (multipart.File, *multipart.FileHeader, error)
```
Повертає дані файлу форми. Використовується із формою ``multipart/form-data``.

__ValidateCsrfToken__
```
ValidateCsrfToken() error
```
Метод який проводить валідацію CSRF token. Для цього форма повинна мати поле із назвою ``csrf_token``, крім того дані cookies
також повинні мати поле ``csrf_token``.<br>
Найпростіший спосіб для цього - додати вбудований [middleware](https://github.com/uwine4850/foozy/blob/master/docs/ua/middlewares.md) який автоматично буде додавати поле ``csrf_token`` в дані cookies.
Після цього просто потрібно додати в середину HTML форми змінну ``{{ csrf_token | safe }}`` та запустити даний метод.<br>
Підключення вбудованого middleware для створення токена буде відбуватись наступним чином:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
...
newRouter.SetMiddleware(mddl)
```

__Files__
```
Files(key string) ([]*multipart.FileHeader, bool)
```
Повертає декілька файлів із форми(multiple input).

## Глобальні функції пакета
__FillStructFromForm__
```
FillStructFromForm(frm interfaces.IForm, fill interface{}) error
```
Метод, який заповнює структуру даними з форми.
Структура завжди повинна передаватись як посилання.
Для коректної роботи необхідно для кожного поля структури вказати тег "form". Наприклад, `form:<ім'я поля форми>`.
Поля структури можуть бути тільки двох типів:
* []FormFile - файли форм.
* []string - всі інші дані.

__FrmValueToMap__
```
FrmValueToMap(frm interfaces.IForm) map[string]interface{}
```
Конвертує дані форму у мапу.

__ReplaceFile__
```
ReplaceFile(pathToFile string, w http.ResponseWriter, file multipart.File, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) error
```
Заміняє файл із файлової системи новим файлом.
