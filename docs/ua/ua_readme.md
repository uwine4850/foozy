[UA](https://github.com/uwine4850/foozy/blob/master/docs/ua/ua_readme.md) | [EN](https://github.com/uwine4850/foozy)<br>
__foozy__ — це легкий та гнучкий веб-фреймворк. В основі проекту лежать модулі http.ServeMux та http.Server. Також,
по можливості, модулі залежать від інтерфейсів, тому більшість з них відкриті для змін.

Модулі які містить фреймворк: <br>
* [builtin](https://github.com/uwine4850/foozy/blob/master/docs/ua/builtin/builtin.md) — вбудований готовий функціонал, наприклад, аутентифікація. Використовувати не обов'язково.
* [database](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/database.md) — інтерфейс для роботи з базою даних mysql.
  * [dbutils](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/dbutils/dbutils.md) — допоміжний функціонал для використання пакету database.
  * [dbmapper](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/dbmapper/dbmapper.md) — записує дані у вибраний об'єкт.
  * [sync_queres](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/sync_queries.md) — синхроні запити до бази даних.
  * [async_queres](https://github.com/uwine4850/foozy/blob/master/docs/ua/database/async_queries.md) — асинхроні запити до бази даних.
* [interfaces](https://github.com/uwine4850/foozy/blob/master/docs/ua/interfaces/interfaces.md) — усі golang інтерфейси які використовуються у проекті.
* [router](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/router.md) — це найголовніший модуль, з допомогою його функціоналу реалізується маршрутизація проекту та багато іншого.
  * [manager](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/manager/manager.md) — пакет для управління обробниками.
  * [websocket](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/websocket.md) — пакет для взаємодії з вебсокетами.
  * [form](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/form/form.md) — робота з HTML формами.
	* [fill_struct](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/form/fill_struct.md) — різноманітні маніпуляції з даними форми.
  * [middlewares](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/middlewares/middlewares.md) — модуль для створення проміжного ПО.
  * [object](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/object/object.md) — пакет для більш простого відображення шаблонів.
  * [mic](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/mic/mic.md) - пакет відповідає за функціональність мікросервісів.
  * [tmlengine](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/tmlengine/tmlengine.md) — шаблонізатор проекту. Використовується бібліотека pongo2.
* [server](https://github.com/uwine4850/foozy/blob/master/docs/ua/server/server.md) — надбудова над http.Server для простішого використання та роботи з модулем router.
  * [livereload](https://github.com/uwine4850/foozy/blob/master/docs/ua/server/livereload/livereload.md) — модуль, який можна використати для перезавантаження проекта після оновлення файлів.
* [utils](https://github.com/uwine4850/foozy/blob/master/docs/ua/utils/utils.md) — загальний допоміжний функціонал, наприклад, генерація CSRF токена.

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
Метод ``NewRouter(manager interfaces.IManager) *Router`` потребує менеджер для роботи, тому код буде виглядати так:
```
newManager := router.NewManager()
newRouter := router.NewRouter(newManager)
```
В свою чергу менеджер для роботи потребує структуру рендера ``NewManager(render interfaces.IRender) *Manager``.
Потрібно її додати:
```
render, err := tmlengine.NewRender()
if err != nil {
    panic(err)
}
newManager := router.NewManager(render)
newRouter := router.NewRouter(newManager)
```
Далі потрібно задати маршрути для роботи. Наприклад, щоб перейти на сторінку __/home__ потрібно зробити такий обробник:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    manager.Render().SetTemplatePath("templates/home.html")
	if err := manager.Render().RenderTemplate(w, r); err != nil {
	    panic(err)
    }
    return func() {}
})
```
Цей код запускає обробник по адресі __/home__. Коли користувач переходить по ній він отримає HTML шаблон по адресі
__templates/home.html__, цей шаблон встановлюється за допомогою ``manager.Render().SetTemplatePath("templates/home.html")``, далі
він відображається з допомогою ``manager.Render().RenderTemplate(w, r)``.<br>
Важливо зазначити, що не обов'язково використовувати шаблонізатор(також можна його змінити), можна використовувати
влаштований метод ``w.Write()`` або інші. Наприклад, є можливість вивести на сторінку дані в форматі JSON, наприклад:
```
newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    values := map[string]string{"key1": "val1"}
	if err := manager.Render().RenderJson(values, w); err != nil {
		panic(err)
	}
	return func() {}
})
```
Для кожної веб сторінки крім HTML потрібен CSS та JavaScript. Його можна додати з допомогою наступного коду:
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
serv := server.NewServer(":8000", newRouter)
err = serv.Start()
if err != nil && !errors.Is(http.ErrServerClosed, err) {
	panic(err)
}
```
Цей код означає, що сервер буде запущено на локальному хості та буде знаходитись на порту 8000. Для відображення сторінок
буде використовуватися маршрутизатор із змінної ``newRouter``.<br>
Повний код міні-проекту наведено нижче.
```
package main

import (
	"errors"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
)

func main() {
	render, err := tmlengine.NewRender()
	if err != nil {
		panic(err)
	}
	newManager := manager.NewManager(render)
	newRouter := router.NewRouter(newManager)
    newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) 
    func() {
        manager.Render().SetTemplatePath("templates/home.html")
        if err := manager.Render().RenderTemplate(w, r); err != nil {
            panic(err)
        }
        return func() {}
    })
	serv := server.NewServer(":8000", newRouter)
	err = serv.Start()
	if err != nil && !errors.Is(http.ErrServerClosed, err) {
		panic(err)
	}
}
```