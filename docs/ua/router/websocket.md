## Package websocket
Інтерфейс веб-сокета реалізований з допомогою бібліотеки __github.com/gorilla/websocket__. 
В пакеті ``router`` знаходиться глобальна змінна ``Upgrader`` яка потрібна 
роботи веб-сокета.

__OnConnect__
```
OnConnect(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
```
Функція яка запускається під час підлючення до клієнта.

__Close__
```
Close() error
```
Закриття підключення.

__OnClientClose__
```
OnClientClose(fn func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn))
```
Функція яка буде виконана коли клієнт закриє з'єднання.

__OnMessage__
```
OnMessage(fn func(messageType int, msgData []byte, conn *websocket.Conn))
```
Коли сокет отримає повідомлення виконається функція ``fn``.

__SendMessage__
```
SendMessage(messageType int, msg []byte, conn *websocket.Conn) error
```
Відправлення повідомлення клієнту.

__ReceiveMessages__
```
ReceiveMessages(w http.ResponseWriter, r *http.Request) error
```
Метод який запускає приймання повідомлень. Цей метод повинен бути запущений.