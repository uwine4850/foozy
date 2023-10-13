## Router package
Модуль router — це важливий модуль, адже на ньому будується весь веб-додаток. Даний модуль ділиться на декілька частин, 
а саме:
* Роутер
* Менеджер
* Вебсокети

## Router
Роутер відповідає за маршрутизацію та роботу із обробниками. Нижче показані методи які доступні для використання.<br>
### Обробники маршрутів
У всіх обробниках є стандартні параметри:
* *pattern* — шлях до обробника. Наприклад, шлях може бути такий ``/home`` і всі до нього подібні. Маршрутизатор також 
підтримує slug параметри, наприклад, ```/post/<id>```. Якщо використати такий шлях і перейти за адресою ``/post/1`` — запуститься 
потрібний обробник в якому буде доступний параметр __id__ у менеджері ось так ``manager.GetSlugParams("id")``. Таких slug 
параметрів може бути багато, головне з різними назвами.
* *fn func(w http.ResponseWriter, r \*http.Request, manager interfaces.IManager)* — функція яка запуститься коли користувач 
перейде на потрібну адресу. ``w http.ResponseWriter`` та ``*http.Request`` це стандартні структури golang. Про ``interfaces.IManager`` 
детальныше написано [тут](#manager).

__Get__
```
Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager))
```
Метод використовується для передачі даних із сервера на веб-сторінку. Наприклад, це може бути html дані, JSON дані та інші.

__Post__
```
Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager))
```
Метод для обробки запита POST. Частіше за все використовується для роботи із формами HTML. Цей метод не має нового функціоналу 
порівняно із методом ``Get``. Але це можна змінити за допомогою [пакету обробника форми](https://github.com/uwine4850/foozy/blob/master/docs/ua/form.md).

__Ws__
```
Ws(pattern string, ws interfaces2.IWebsocket, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager))
```
Обробник запускає веб-сокет по вибраному шляху. До цього обробника можна без проблем під'єднатися з допомогою JavaScript.
Параметр ``interfaces2.IWebsocket`` це інтерфейс структури яка реалізує взаємодію з веб-сокетом, ось стандартна реалізація
``router.NewWebsocket(router.Upgrader)``. Про веб-сокети детальніше написано [тут](#websocket).<br>
Приклад реалізації echo обробника:
```
newRouter.Ws("/ws", router.NewWebsocket(router.Upgrader), func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
    ws := manager.GetWebSocket()
	err := ws.Connect(w, r, func() {
	    fmt.Println("Connect.")
	})
	if err != nil {
	    panic(err)
	}
	ws.OnClientClose(func() {
	    err := ws.Close()
		if err != nil {
		    panic(err)
		}
		fmt.Println("Client close.")
	})
	ws.OnMessage(func(messageType int, msgData []byte) {
	    err := ws.SendMessage(messageType, msgData)
		if err != nil {
		    panic(err)
		}
	})
	err = ws.ReceiveMessages()
	if err != nil {
	    panic(err)
	}
})
```

### Інші методи
__EnableLog__
```
EnableLog(enable bool)
```
Ввімкнути або вимкнути виведення логу в консоль.

__SetMiddleware__
```
SetMiddleware(middleware interfaces.IMiddleware)
```
Встановлює екземпляр інтерфейсу ``interfaces.IMiddleware`` для запуску перед кожним обробником.

__getHandleFunc__
```
getHandleFunc(pattern string, method string, fn func(w http.ResponseWriter, r *http.Request, manager interfaces2.IManager)) http.HandlerFunc
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

## Manager
Даний компонент відповідає за менеджмент обробника. На цей момент інтерфейс виконує наступні функції:
* Передача даних із middlewares у обробник маршруту.
* Замінити стандартний шаблонізатор.
* Рендер шаблона.
* Отримання slug параметрів.
* Доступ до інтерфейсу веб-сокетів.
* Рендер даних у форматі JSON.

### Методи інтерфейсу
__SetUserContext__
```
SetUserContext(key string, value interface{})
```
Встановлює користувацький контекст, наприклад, в проміжному ПО.

__GetUserContext__
```
GetUserContext(key string) (any, bool)
```
Повертає значення користувацького контексту по ключу.

__SetTemplateEngine__
```
SetTemplateEngine(engine interfaces2.ITemplateEngine)
```
Змінює стандартний шаблонізатор.

__RenderTemplate__
```
RenderTemplate(w http.ResponseWriter, r *http.Request) error
```
Відображає шаблон з допомогою шаблонізатора.<br>
__ВАЖЛИВО:__ шаблон потрыбно встановити з допомогою метода ``SetTemplatePath``.

__SetTemplatePath__
```
SetTemplatePath(templatePath string)
```
Встановлює шлях до HTML шаблона.

__SetContext__
```
SetContext(data map[string]interface{})
```
Встановлює контекст для шаблонізатора. В HTML шаблоні це виглядає так ``{{ key }}``.

__SetSlugParams__
```
SetSlugParams(params map[string]string)
```
Встановлює slug параметри. Використовується в роутері.

__GetSlugParams__
```
GetSlugParams(key string) (string, bool)
```
Дає доступ до slug параметрів.

__SetWebsocket__
```
SetWebsocket(websocket interfaces.IWebsocket)
```
Встановлює веб-сокет. Використовується в роутері.

__GetWebSocket__
```
GetWebSocket() interfaces.IWebsocket
```
Надає доступ до інтерфейсу веб-сокета.

__RenderJson__
```
RenderJson(data interface{}, w http.ResponseWriter) error
```
Відображає дані у форматі JSON. Як параметр data може приймати мапу, структуру та інше.

## Websocket
Інтерфейс веб-сокета реалізований з допомогою бібліотеки __github.com/gorilla/websocket__. В пакеті ``router`` є глобальна 
змінна ``Upgrader`` яка потрібна роботи веб-сокета.

### Методи інтерфейсу
__Connect__
```
Connect(w http.ResponseWriter, r *http.Request, fn func()) error
```
Підключення до клієнта(наприклад JavaScript). Параметр ``fn func()`` це функція яка буде виконана під час кожного підключення.

__Close__
```
Close() error
```
Закриття підключення.

__OnClientClose__
```
OnClientClose(fn func())
```
Функція яка буде виконана коли клієнт закриє з'єднання.

__OnMessage__
```
OnMessage(fn func(messageType int, msgData []byte))
```
Коли сокет отримає повідомлення виконається функція ``fn``.

__SendMessage__
```
SendMessage(messageType int, msg []byte) error
```
Відправлення повідомлення клієнту.

__ReceiveMessages__
```
ReceiveMessages() error
```
Метод який запускає приймання повідомлень. Цей метод повинен бути запущений.
