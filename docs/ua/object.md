## Form package
Даний пакет містить в собі функції та інтерфейси для більш зручнішого керування шаблонами.

__type TemplateView struct__

Структура для запуска будь-якого інтерфейсу IView.

* _Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()_ - метод, який запускає виконання 
обробника http запиту.

__type IView interface__

Інтерфейс, який повинен реалізовувати кожен View.

* _Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error)_ - метод 
звертається до бази даних та записує їх у контекст шаблонізатора.<br>
* _GetContext() map[string]interface{}_ - повертає контекст який буде встановлений у шаблонізатор.<br>
* _Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}_ - метод який 
потрібно перевизначити у користувацькій структурі. Важливим моментом є те, що метод __Object__ записує дані у контекст
перед виконанням цього метода, тому для отримання об'єкта із бази даних можна використати метод __GetContext__.<br>
* _Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()_ - для використання метода 
його потрібно перевизначити. З домомогою цього метода можна визначити права доступу для адреси. У випадку блокування доступу
потрібно повернути false та фунцію яку потрібно виконати.<br>
* _OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)_ - метод потрібно перевизначити. 
Даний метод буде виконаний коли під час внутрішнього виконання алгоритмів виникне помилка.<br>

## type ObjView struct

Детальний перегляд конкретного запису із бази даних. Приклад використання:
```
type ProfileView struct {
    object.ObjView
}

func (v *ProfileView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
    fmt.Println(v.GetContext())
    return map[string]interface{}{"id": 50000}
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

## type MultipleObjectView struct

Виконує ті ж самі дії, що і __ObjView__ з однією відмінністю - знаходження записів у базі даних відбувається у декількох 
таблицях. Приклад використання:
```
type ProfileMulView struct {
    object.MultipleObjectView
}

func (v *ProfileMulView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
    return map[string]interface{}{}
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
                        "profile",
                        "auth",
                        "id",
                        "id",
                        User{},
                    },
                    {
                        Name:       "tee",
                        TaleName:   "tee",
                        SlugName:   "tee",
                        AIField:    "id",
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
_MultipleObjects_ - екземпляр структури MultipleObjects.<br>

__MultipleObject__

Структура яка пердставляє дані про конкретний запис у базі даних.

_Name_ - назва з допомогою якого можна звернутись до запусу із бази даних.<br>
_TaleName_ - назва таблиці.<br>
_SlugName_ - назва slug для отримання його значення.<br>
_AIField_ - назва колонки по якій потрібно шукати значення у таблиці з допомогою slug.
_FillStruct_ - структура яка описує таблицю.<br>

## type AllView struct

Виводить усі дані із таблиці. Приклад використання:
```
type ProjectView struct {
    object.AllView
}

func (v *ProjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
    return map[string]interface{}{}
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
