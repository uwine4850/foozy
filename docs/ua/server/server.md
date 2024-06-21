## Package server
Цей пакет містить інтерфейс для роботи з ``http.Server`` та ``interfaces.IRouter``.<br>
Приклад роботи:
```
server := fserer.NewServer(":8000", newRouter)
err = server.Start()
if err != nil {
    panic(err)
}
```
Для зупинки можна користуватися ``ctrl + c``.

## Методи
__Start__
```
Start() error
```
Запуск сервера.

__GetServ__
```
GetServ() *http.Server
```
Повертає екземпляр ``*http.Server``.

__Stop__
```
Stop() error
```
Зупинка сервера.

__type MicServer struct__

Структура для запуску мікросервіса.

* _Start() error_ - запуск сервера.

__FoozyAndMic__
```
FoozyAndMic(fserver *Server, micServer *MicServer, onError func(err error))
```
Запуск одразу звичайного сервера та сервера мікросервіса.
