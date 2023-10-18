[UA](https://github.com/uwine4850/foozy/blob/master/docs/ua/ua_readme.md) | [EN](https://github.com/uwine4850/foozy)<br>
__foozy__ — це легкий та гнучкий веб-фреймворк. В основі проекту лежать модулі http.ServeMux та http.Server. Також,
по можливості, модулі залежать від інтерфейсів, тому більшість з них відкриті для змін.

Модулі які містить фреймворк: <br>
* [builtin](https://github.com/uwine4850/foozy/blob/master/docs/ua/builtin.md) — вбудований готовий функціонал, наприклад, аутентифікація. Використовувати не обов'язково.
* [database](https://github.com/uwine4850/foozy/blob/master/docs/ua/database.md) — інтерфейс для роботи з базою даних mysql.
* interfaces — усі golang інтерфейси які використовуються у проекті.
* [livereload](https://github.com/uwine4850/foozy/blob/master/docs/ua/livereload.md) — модуль, який можна використати для перезавантаження проекта після оновлення файлів.
* [middlewares](https://github.com/uwine4850/foozy/blob/master/docs/ua/middlewares.md) — модуль для створення проміжного ПО.
* [router](https://github.com/uwine4850/foozy/blob/master/docs/ua/router.md) — це найголовніший модуль, з допомогою його функціоналу реалізується маршрутизація проекту та багато іншого.
* [form](https://github.com/uwine4850/foozy/blob/master/docs/ua/form.md) — робота з HTML формами.
* [server](https://github.com/uwine4850/foozy/blob/master/docs/en/server.md) — надбудова над http.Server для простішого використання та роботи з модулем router.
* tmlengine — шаблонізатор проекту. Використовується бібліотека pongo2.
* utils — загальний допоміжний функціонал, наприклад, генерація CSRF токена.

## Початок роботи

### Встановлення
```
go get github.com/uwine4850/foozy
```

### Базове використання
Для початку потрібно використати роутер ``router.Router``, наприклад:
```
newRouter := router.NewRouter()
```
Метод ``NewRouter(manager interfaces2.IManager) *Router`` потребує менеджер для роботи, тому код буде виглядати так:
```
newManager := router.NewManager()
newRouter := router.NewRouter(newManager)
```
В свою чергу менеджер для роботи потребує шаблонізатор ``NewManager(engine interfaces2.ITemplateEngine) *Manager``.
Потрібно його додати:
```
newTmplEngine, err := tmlengine.NewTemplateEngine()
if err != nil {
    panic(err)
}
newManager := router.NewManager(newTmplEngine)
newRouter := router.NewRouter(newManager)
```
Далі потрібно задати маршрути для роботи. Наприклад, щоб перейти на сторінку __/home__ потрібно зробити такий обробник:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    manager.SetTemplatePath("templates/home.html")
    err := manager.RenderTemplate(w, r)
    if err != nil {
        panic(err)
    }
})
```
Цей код запускає обробник по адресі __/home__. Коли користувач переходить по ній він отримає HTML шаблон по адресі
__templates/home.html__, цей шаблон встановлюється за допомогою ``manager.SetTemplatePath("templates/home.html")``, далі
він відображається з допомогою ``manager.RenderTemplate(w, r)``.<br>
Важливо зазначити, що не обов'язково використовувати шаблонізатор(також можна його змінити), можна використовувати
влаштований метод ``w.Write()`` або інші. Наприклад, є можливість вивести на сторінку дані в форматі JSON, наприклад:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    values := map[string]string{"key1": "val1"}
    err = manager.RenderJson(values, w)
    if err != nil {
        panic(err)
    }
}
```
Для кожного сайту крім HTML потрібен CSS та JavaScript. Його можна додати з допомогою наступного коду:
```
newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
```
Даний код розраховує на те, що усі файли будуть знаходитись у директорії static. В HTML підключення файлів буде виглядати так:
```
<link rel="stylesheet" href="/static/css/style.css">
...
<img src="/static/img/image.png">
```
Важливо зазначити, що шлях завжди повинен починатися з символа ``/``.<br>
Отже, тепер коли є базовий обробник потрібно запустити сервер наступним чином.
```
_server := server.NewServer(":8000", newRouter)
err = _server.Start()
if err != nil {
    panic(err)
}
```
Цей код означає, що сервер буде запущено на локальному хості та буде знаходитись на порту 8000. Для відображення сторінок
буде використовуватися маршрутизатор із змінної ``newRouter``.<br>
Повний код міні-проекту наведено нижче.
```
package main

import (
    "github.com/uwine4850/foozy/pkg/interfaces"
    "github.com/uwine4850/foozy/pkg/router"
    "github.com/uwine4850/foozy/pkg/server"
    "github.com/uwine4850/foozy/pkg/tmlengine"
    "net/http"
)

func main() {
    newTmplEngine, err := tmlengine.NewTemplateEngine()
    if err != nil {
        panic(err)
    }
    newManager := router.NewManager(newTmplEngine)
    newRouter := router.NewRouter(newManager)
    newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
        manager.SetTemplatePath("templates/home.html")
        err := manager.RenderTemplate(w, r)
        if err != nil {
            panic(err)
        }
    })
    newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    
    _server := server.NewServer(":8000", newRouter)
    err = _server.Start()
    if err != nil {
        panic(err)
    }
}
```