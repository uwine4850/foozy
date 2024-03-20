## Middlewares package
Даний пакет містить всі потрібні інструменти для роботи з проміжним ПО.

## Методи які використовує інтерфейс IMiddleware
Методи які створюють обробник, мають один загальний параметр ``fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData)`` - 
цей параметр ідентичний із параметром який описується в пакеті [router](https://github.com/uwine4850/foozy/blob/master/docs/ua/router.md).<br>
Єдина відмінність полягає у тому, що дані цього параметру спочатку потрапляють у middlewares, а потім у обробник роутера.

__HandlerMddl__
```
HandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData))
```
Метод створює проміжне ПО яке буде виконуватися синхронно. Параметр ``id`` - це порядковий номер виконання middleware, він 
повинен бути унікальним.

__AsyncHandlerMddl__
```
AsyncHandlerMddl(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData))
```
За допомогою цього метода також можна створити проміжне ПО, але воно буде запускатися асинхронно. Відповідно порядкових номерів не існує.

__RunMddl__
```
RunMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error
```
Запуск синхронних middleware.

__RunAsyncMddl__
```
RunAsyncMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
Запуск асинхронних middleware.<br>
__ВАЖЛИВО__: усі проміжні ПО запускаються асинхронно, тому потрібно дочекатися їх виконання за допомогою метода ``WaitAsyncMddl``.

__WaitAsyncMddl__
```
WaitAsyncMddl()
```
Чекати виконання усіх асинхронних проміжних ПО(якщо вони є).

### Методи пакету
__SetMddlError__
```
SetMddlError(mddlErr error, manager interfaces.IManagerData)
```
Зберігає помилку яка виникла у проміжному ПО.

__GetMddlError__
```
GetMddlError(manager interfaces.IManagerData) (error, error)
```
Повертає помилку яка була встановлена за допомогою __SetMddlError__.

__SkipNextPage__
```
SkipNextPage(manager interfaces.IManagerData)
```
Пропускає рендер наступного веб-запиту. Обробник не запускається навіть на початковій стадії.

__IsSkipNextPage__
```
IsSkipNextPage(manager interfaces.IManagerData) bool
```
Перевіряє чи потрібно пропускати наступну сторінку. У стандартній реалізації використовується у роутері.

__SkipNextPageAndRedirect__
```
SkipNextPageAndRedirect(manager interfaces.IManagerData, w http.ResponseWriter, r *http.Request, path string)
```
Пропускає рендер наступної сторінки та перенаправляє на іншу.
