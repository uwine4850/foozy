## Package server
This package contains an interface to work with ``http.Server'' and ``interfaces.IRouter''.<br
Example of work:
```
server := fserer.NewServer(":8000", newRouter)
err = server.Start()
if err != nil {
    panic(err)
}
```
To stop, you can use 'ctrl + c'.

## Методи
__Start__
```
Start() error
```
Start the server.

__GetServ__
```
GetServ() *http.Server
```
Returns an instance of ``*http.Server``.

__Stop__
```
Stop() error
```
Stopping the server.
