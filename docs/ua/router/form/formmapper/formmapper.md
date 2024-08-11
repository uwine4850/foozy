## package formmapper
Заповнює структуру даними із форми.

Ви можете побачити, як працює пакет у цих [тестах](https://github.com/uwine4850/foozy/tree/master/tests/formtest/formmapping_test).

### type Mapper struct

`Form` — посилання на уже парсену форму.<br>
`Output` — посилання на структуру, або структуру у `*reflect.Value`.<br>
`NilIfNotExist` — якщо у формі не знайдені відповідні поля, і вони є у цьому списку, поля структури будуть мати значення nil.

Приклад налаштування структури:
```
type Fill struct {
	Field1   []string        `form:"f1"   empty:"1"`
	Field2   []string        `form:"f2"   empty:"-err"`
	File     []form.FormFile `form:"file" empty:"-err"`
}
```
Кожне поле структури повине бути типом `[]string` або `[]form.FormFile`. Це потрібно для того, 
щоб по одній назві(ключу) input можна було передати декілька зачень.<br>

Тег `form` обов'язково потрібен. Він відповідає за назву input, який буде записаний у поле структури.<br>

Тег `empty` буде застосований лише для порожніх значень. Тут важливо зазначити, що це застосовується лише 
для значення зрізу, тобто, якщо по ключу є два значення, а порожнє 
тільки одне, то операція буде проведена тільки із порожнім значенням. Даний тег має 
декілька опцій:
*  -err — виводить помилку, якщо хочаб один із індексів зрізу порожній. 
Тип `[]form.FormFile` може мати тільки одну опцію — `-err`.
*  просто текст — замініє дані пустого значення за його індексом.

Тег `ext` розширення фалів, які можуть бути у полі. Наприклад `ext:".jpg .png"`.

__Fill__
```
Fill() error
```
Заповнює структуру значеннями із форми.

### type OrderedForm struct

Структура, яка впорядковує форму для подальшого більш зручного використання. Усі поля впорядковані по порядку їхньому порядку у формі.

* _Add(name string, value interface{})_ - додає нове поле форми у структуру.<br>
* _GetByName(name string) (OrderedFormValue, bool)_ - повертає поле форми по його назві.<br>
* _GetAll() []OrderedFormValue_ - повертає усі поля форми.<br>
  
### Інші функції пакету

__FillStructFromForm__
```
FillStructFromForm(frm *Form, fillStruct interface{}, nilIfNotExist []string) error
```
Метод, який заповнює структуру даними з форми.
Структура завжди повинна передаватись як посилання.
Для коректної роботи необхідно для кожного поля структури вказати тег "form". Наприклад, `form:<ім'я поля форми>`. Також підтримує тег `empty`, який описано вище.
* frm *Form - екземпляр форми.
* fillStruct interface{} - посилання на об'єкт, який потрібно заповнити.
* nilIfNotExist - поля які не знайшлись у формі будуть nil.

__FrmValueToOrderedForm__
```
FrmValueToOrderedForm(frm IFormGetEnctypeData) *OrderedForm
```
Заповнює дані форми у структуру *OrderedForm*.

__FieldsNotEmpty__
```
FieldsNotEmpty(fillStruct interface{}, fieldsName []string) error
```
Перевіряє чи не порожні вибрані поля структури.
Оптимізовано для роботи, навіть якщо FillableFormStruct містить структуру з типом *reflect.Value.

__FieldsName__
```
FieldsName(fillStruct interface{}, exclude []string) ([]string, error)
```
Повертає назви полів структури.

__CheckExtension__
```
CheckExtension(fillStruct interface{}) error
```
Перевіряє чи розширення файлів форми відповідає очікуваними. Для правильної роботи потрібно додати до кожного поля типу 
FormFile тег *ext* і розширення які очікуються. Наприклад, `ext:".jpeg .png"`.

__FillReflectValueFromForm__
```
FillReflectValueFromForm(frm *Form, fillValue *reflect.Value, nilIfNotExist []string) error
```
Заповнює структуру даними з форми.
Функція працює і робить все так само, як функція `FillStructFromForm`.
Єдина відмінність полягає в тому, що ця функція приймає дані у форматі `*reflect.Value`.