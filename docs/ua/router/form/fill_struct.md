## Fill struct

Тести для форми [тут](https://github.com/uwine4850/foozy/tree/master/tests/formtest).

__type FillableFormStruct struct__

Структура FillableFormStruct призначена для більш зручного доступу до заповнюваної структури.
Структура, яку потрібно заповнити, передається вказівником.

* _GetStruct() interface{}_ - Повертає заповнену структуру.<br>
* _SetDefaultValue(val func(name string) interface{})_ - встановлює стандартну функцію.<br>
* _GetOrDef(name string, index int) interface_ - повертає значення структури або стандартну функцію якщо значення структури відсутнє.<br>

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
