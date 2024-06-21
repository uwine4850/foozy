## IManagerWebsocket
Надає доступ до структури [вебсокета](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/websocket.md).

__CurrentWebsocket__
```
CurrentWebsocket() IWebsocket
```
Повертає структуру, яка реалізує інтерфейс IWebsocket.

__SetWebsocket__
```
SetWebsocket(websocket IWebsocket)
```
Встановлює структуру, яка реалізує інтерфейс IWebsocket. У стандартній реалізації 
використовується для внутрішніх потреб роутера.