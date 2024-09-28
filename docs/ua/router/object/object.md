## Object package
Даний пакет містить в собі функції та інтерфейси для більш зручнішого керування шаблонами.
Для початку роботи із object потрібно створити нову структуру із вбудувати у неї вибраний object. Важливо зазначити, що може бути вбудований 
лише один object.

### type IView interface
Інтерфейс, який повинен реалізовувати кожен View.

* _Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error)_ - метод 
звертається до бази даних та записує їх у контекст шаблонізатора.<br>
* _Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error)_ - метод який 
потрібно перевизначити у користувацькій структурі. Важливим моментом є те, що метод __Object__ записує дані у контекст
перед виконанням цього метода, тому для отримання даних які встановлені у Object потрібно використати метод 
manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT). У цьому методі також доступне активне підключення до бази даних,
його можна отримати методо GetDB().<br>
* _Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()_ - для використання метода 
його потрібно перевизначити. З домомогою цього метода можна визначити права доступу для адреси. У випадку блокування доступу
потрібно повернути false та фунцію яку потрібно виконати.<br>
* _OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)_ - метод потрібно перевизначити. 
Даний метод буде виконаний коли під час внутрішнього виконання алгоритмів виникне помилка.<br>
* _ObjectsName() []string_ - повертає назви об'єктів, або одного об'єкта.

__GetObjectContext__
```
GetObjectContext(manager interfaces.IManager) (ObjectContext, error)
```
Повертає з менеджера ObjectContext.
Важливо розуміти, що цей метод може бути використаний тільки коли метод IView.Object завершив роботу, наприклад, IView.Context.

## Відображення View як HTML.

__type TemplateView struct__

Структура для запуска будь-якого інтерфейсу IView.

* _Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()_ - метод, який запускає виконання 
обробника http запиту.

## Відображення View як JSON.
Кожна з цих структур має метод Call, який використовується для її запуску.
Також, кожна структура відображає дані тільки в форматі JSON.
Якщо не передавати дані у поле `message` - JSON буде відображатись без валідації DTO та у форматі `TemplateView`.

__type JsonObjectTemplateView struct__<br>
Структура яка відображає ObjectView.

__type JsonMultipleObjectTemplateView struct__<br>
Структура яка відображає MultipleObjectView.

__type JsonAllTemplateView struct__<br>
Структура яка відображає AllView.

## Відображення об'єктів із БД у різних форматах.

### type ObjView struct
Детальний перегляд конкретного запису із бази даних. Приклад використання:
```
type ProfileView struct {
    object.ObjView
}

func (v *ProfileView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
    fmt.Println(v.GetContext())
    return ObjectContext{"id": 50000}, nil
}

func Init() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
    view := object.TemplateView{
        TemplatePath: "project/templates/profile.html",
        View: &ProfileView{
            object.ObjView{
                Name:       "profile",
                DB:         db,
                TableName:  "auth",
                FillStruct: User{},
                Slug:       "id",
            },
        },
    }
    return view.Call
}
```
_Name_ - назва з допомогою якого можна звернутись до запусу із бази даних.<br>
_DB_ - екземпляр бази даних.<br>
_TableName_ - назва таблиці.<br>
_FillStruct_ - структура яка описує таблицю.<br>
_Slug_ - slug значення по якому потрібно знайти значення в талиці.<br>

### type MultipleObjectView struct
Виконує ті ж самі дії, що і __ObjView__ з однією відмінністю - знаходження записів у базі даних відбувається у декількох 
таблицях. Приклад використання:
```
type ProfileMulView struct {
    object.MultipleObjectView
}

func (v *ProfileMulView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
    return ObjectContext{}, nil
}

func Init() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
    view := object.TemplateView{
        TemplatePath: "project/templates/profile.html",
        View: &ProfileMulView{
            object.MultipleObjectView{
                DB: db,
                MultipleObjects: []object.MultipleObject{
                    {
                        Name:       "profile",
                        TaleName:   "auth",
                        SlugName:   "id",
                        SlugField:  "id",
                        FillStruct: User{},
                    },
                    {
                        Name:       "tee",
                        TaleName:   "tee",
                        SlugName:   "tee",
                        SlugField:    "id",
                        FillStruct: Tee{},
                    },
                },
            },
        },
    }
    return view.Call
}
```
_DB_ - екземпляр бази даних.<br>
_MultipleObjects_ - екземпляр структури `MultipleObjects`.<br>

`MultipleObject`<br>
Структура яка пердставляє дані про конкретний запис у базі даних.

_Name_ - назва з допомогою якого можна звернутись до запусу із бази даних.<br>
_TaleName_ - назва таблиці.<br>
_SlugName_ - назва slug для отримання його значення.<br>
_SlugField_ - назва колонки по якій потрібно шукати значення у таблиці з допомогою slug.
_FillStruct_ - структура яка описує таблицю.<br>

### type AllView struct
Виводить усі дані із таблиці. Приклад використання:
```
type ProjectView struct {
    object.AllView
}

func (v *ProjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
    return ObjectContext{}, nil
}

func Init() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
    view := object.TemplateView{
        TemplatePath: "project/templates/profile.html",
        View: &ProjectView{
            object.AllView{
                Name:       "data",
                DB:         db,
                TableName:  "project",
                FillStruct: Project{},
            },
        },
    }
    return view.Call
}
```

### type FormView struct
Даний об'єкт трохи відрізняється від інших об'єктів. Відмінність полягає в тому, що головне призначення цього 
об'єкту - прочитати та обробити дані форми. Форма обробляється у методі __Object__ і потім передається у контекст.
Отримати значення у методі `Context` можна як завжди з допомогою `manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT)`.
Також при потребі можна відобразити HTML стрінку, але це не обов'язково. Щоб не відображати сторінку потрібно викликати 
метод `object.TemplateView.SkipRender()`.
Для перенаправлення на іншу сторінку можна викликати `http.Redirect` прямо у методі `Context`.

Важливо зауважити, що параметр `NotNilFormFields` універсальний. Якщо у нього передати першим елементом знак "*", усі поля
структури будуть перевірятися на пустоту. Якщо після цього знаку вказати ще поля структури, вони будуть будуть виключені
з перевірки на пустоту. Також можна просто передавати поля які потрібно перевірити без знаку "*".

Якщо параметр `ValidateCSRF` дорівнює `true` - буде відбуватися перевірка CSRF токена.

__FormInterface__
```
FormInterface(manager interfaces.IManagerOneTimeData) (interface{}, error)
```
Повертає заповнену форму у форматі інтерфейсу.

Приклад використання:
```
type ObjectForm struct {
	Text []string        `form:"text"`
	File []form.FormFile `form:"file"`
}

type MyFormView struct {
	object.FormView
}

func (v *MyFormView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	filledFormInterface, err := v.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	filledForm := filledFormInterface.(ObjectForm)
	if filledForm.Text[0] != "field" {
		return nil, errors.New("FormView unexpected text field value")
	}
	if filledForm.File[0].Header.Filename != "x.png" {
		return nil, errors.New("FormView unexpected file field value")
	}
	return object.ObjectContext{}, nil
}

func (v *MyFormView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func MyFormViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	tv := object.TemplateView{
		TemplatePath: "",
		View: &MyFormView{
			object.FormView{
				FormStruct:       ObjectForm{},
				NotNilFormFields: []string{"*"},
				NilIfNotExist:    []string{},
			},
		},
	}
	tv.SkipRender()
	return tv.Call
}
```