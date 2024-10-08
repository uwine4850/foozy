## Form package
Пакет form містить весь доступний функціонал(на даний момент) для роботи з формами.

Тести для форми [тут](https://github.com/uwine4850/foozy/tree/master/tests/formtest).

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

__Files__
```
Files(key string) ([]*multipart.FileHeader, bool)
```
Повертає декілька файлів із форми(multiple input).

## Глобальні функції пакета

__SaveFile__
```
SaveFile(w http.ResponseWriter, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) error
```
Зберігає файл у вибраному місці.
* fileHeader *multipart.FileHeader - інформація про файл.
* pathToDir string - шлях до директорії для збереження файлу.
* buildPath *string - посилання на зміну string. У зміну записується повний шлях до збереженого файлу.

__ReplaceFile__
```
ReplaceFile(pathToFile string, w http.ResponseWriter, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) error
```
Заміняє вже існуючий файл іншим.

__SendApplicationForm__
```
SendApplicationForm(url string, values map[string]string) (*http.Response, error)
```
Надсилає POST запит(application/x-www-form-urlencoded) із даними по вибраному url.

__SendMultipartForm__
```
SendMultipartForm(url string, values map[string]string, files map[string][]string) (*http.Response, error)
```
Надсилає POST запит(multipart/form-data) із даними по вибраному url.
