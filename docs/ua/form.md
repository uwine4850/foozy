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


### Fill struct

__type FillableFormStruct struct__

Структура FillableFormStruct призначена для більш зручного доступу до заповнюваної структури.
Структура, яку потрібно заповнити, передається вказівником.

* _GetStruct() interface{}_ - Повертає заповнену структуру.<br>
* _SetDefaultValue(val func(name string) string)_ - встановлює стандартну функцію.<br>
* _GetOrDef(name string, index int) string_ - повертає значення структури або стандартну функцію якщо значення структури відсутнє.<br>

__FillStructFromForm__
```
FillStructFromForm(frm *Form, fillableStruct *FillableFormStruct, nilIfNotExist []string) error
```
Метод, який заповнює структуру даними з форми.
Структура завжди повинна передаватись як посилання.
Для коректної роботи необхідно для кожного поля структури вказати тег "form". Наприклад, `form:<ім'я поля форми>`.
* frm *Form - екземпляр форми.
* fillableStruct *FillableFormStruct - екземпляр FillableFormStruct.
* nilIfNotExist - поля які не знайшлись у формі будуть nil.

__type OrderedForm struct__

Структура, яка впорядковує форму для подальшого більш зручного використання. Усі поля впорядковані по порядку їхньому порядку у формі.

* _Add(name string, value interface{})_ - додає нове поле форми у структуру.<br>
* _GetByName(name string) (OrderedFormValue, bool)_ - повертає поле форми по його назві.<br>
* _GetAll() []OrderedFormValue_ - повертає усі поля форми.<br>

__FrmValueToOrderedForm__
```
FrmValueToOrderedForm(frm IFormGetEnctypeData) *OrderedForm
```
Заповнює дані форми у структуру *OrderedForm*.

__FieldsNotEmpty__
```
FieldsNotEmpty(fillableStruct *FillableFormStruct, fieldsName []string) error
```
Перевіряє чи не порожні вибрані поля структури.

__FieldsName__
```
FieldsName(fillForm *FillableFormStruct, exclude []string) ([]string, error)
```
Повертає назви полів структури.

__CheckExtension__
```
CheckExtension(fillForm *FillableFormStruct) error
```
Перевіряє чи розширення файлів форми відповідає очікуваними. Для правильної роботи потрібно додати до кожного поля типу 
FormFile тег *ext* і розширення які очікуються. Наприклад, `ext:".jpeg, .png"`.
