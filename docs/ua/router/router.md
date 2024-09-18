## Router package

## Router
Роутер відповідає за маршрутизацію та роботу із обробниками. Нижче показані методи які доступні для використання.<br>

Тести для роутера [тут](https://github.com/uwine4850/foozy/tree/master/tests/routing).

### Обробники маршрутів
У всіх обробниках є стандартні параметри:
* *pattern* — шлях до обробника. Наприклад, шлях може бути такий ``/home`` і всі до нього подібні. Маршрутизатор також 
підтримує slug параметри, наприклад, ```/post/<id>```. Якщо використати такий шлях і перейти за адресою ``/post/1`` — запуститься 
потрібний обробник в якому буде доступний параметр __id__ у менеджері ось так ``manager.GetSlugParams("id")``. Таких slug 
параметрів може бути багато, головне з різними назвами.
* *fn func(w http.ResponseWriter, r \*http.Request, manager interfaces.IManager)* — функція яка запуститься коли користувач 
перейде на потрібну адресу. ``w http.ResponseWriter`` та ``*http.Request`` це стандартні структури golang. Про ``interfaces.IManager`` 
детальныше написано [тут](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/manager/manager.md).
* Кожен обробник повертає ``func()`` - це функція, яка виконується після завершення роботи самого обробника.

Для одного маршруту може бути застосовані декілька обробників методів. 
Головне щоб методи не повторювались, тобто тільки один Get, Post, Delete і тд.

__Get__
```
Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
Метод використовується для передачі даних із сервера на веб-сторінку. Наприклад, це може бути html дані, JSON дані та інші.

__Post__
```
Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
Метод для обробки запита POST. Частіше за все використовується для роботи із формами HTML. Цей метод не має нового функціоналу 
порівняно із методом ``Get``. Але це можна змінити за допомогою [пакету обробника форми](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/form/form.md).

__Ws__
```
Ws(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) func()
```
Обробник приймає веб-сокет по вибраному шляху. Про веб-сокети детальніше написано [тут](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/websocket.md).<br>
Приклад реалізації echo обробника:
```
newRouter.Ws("/ws", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	ws := router.NewWebsocket(router.Upgrader)
	ws.OnConnect(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		fmt.Println("Connect.")
	})
	ws.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn websocket.Conn) {
		err := ws.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Client close.")
	})
	ws.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		err := ws.SendMessage(messageType, msgData, conn)
		if err != nil {
			panic(err)
		}
	})
	err = ws.ReceiveMessages(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
})
```

__Put__
```
Put(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func())
```
Метод для обробки запита PUT.

__Delete__
```
Delete(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
```
Метод для обробки запита Delete.

__Options__
```
Options(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
```
Метод для обробки запита Options.

__InternalError__
```
InternalError(fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error))
```
InternalError встановлює функцію, яка використовуватиметься під час обробки внутрішніх помилок.

### Інші методи

__RegisterAll__
```
RegisterAll()
```
Реєструє всі обробники.

__SetMiddleware__
```
SetMiddleware(middleware interfaces.IMiddleware)
```
Встановлює екземпляр інтерфейсу ``interfaces.IMiddleware`` для запуску перед кожним обробником.

__getHandleFunc__
```
getHandleFunc(pattern string, method string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) http.HandlerFunc
```
Приватний метод який запускається перед кожним обробником. Він запускає різноманітні валідації, проміжне ПО та інше.
Цей метод обгортає функцію ``func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)`` яка передана із 
обробника Post, Get або Ws.

__validateMethod__
```
validateMethod(method string) bool
```
Валідація http методів які зараз використовується. Тобто, слідкує щоб метод Get оброблював метод http GET, а не, наприклад, POST.

#### Методи, які не належать інтерфейсу, але належать пакету
Ці методи є глобальними, знаходяться в пакеті router, але можуть використовуватися будь-де.<br>

__ValidateRootUrl__
```
ValidateRootUrl(w http.ResponseWriter, r *http.Request) bool
```
Якщо паттерн шляху дорівнює __/__, то обробник буде приймати __усі шляхи__. Щоб запобігти цьому потрібно використати цей 
метод. Тепер якщо шлях не буде знайдено буде виводитись помилка 404, а не обробник __/__. Приклад:
```
newRouter.Get("/", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    if !router.ValidateRootUrl(w, r) {
	    return
	}
})
```
__ParseSlugIndex__
```
ParseSlugIndex(path []string) map[int]bool
```
Функція ділить url на частини по символу __/__. Далі він по числовому порядку додає кожну частину в карту як ключ, а значення 
ключа залужить від того чи є він slug параметром, якщо це так ключ дорівнює true, якщо ні — false.

__HandleSlugUrls__
```
HandleSlugUrls(parseUrl map[int]bool, slugUrl []string, url []string) (string, map[string]string)
```
``slugUrl []string`` це параметр pattern який розділений по символу __/__ та записаний у slice.
``url []string`` це справжній url який розділений по символу __/__ та записаний у slice.

Функція опрацьовує url та виводить рядок(str) як url який зроблений із паттерна та параметри slug якщо вони є.<br>
За допомогою ``parseUrl map[int]bool`` знаходяться частини які є slug параметрами, та які потрібно замінити. Дані для зміни 
беруться із справжнього url по їх числовій позиції.
