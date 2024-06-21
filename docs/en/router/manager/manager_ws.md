## IManagerWebsocket
Provides access to the [websocket structure](https://github.com/uwine4850/foozy/blob/master/docs/en/router/websocket.md).

__CurrentWebsocket__
```
CurrentWebsocket() IWebsocket
```
Returns a structure that implements the IWebsocket interface.

__SetWebsocket__
```
SetWebsocket(websocket IWebsocket)
```
Sets a structure that implements the IWebsocket interface. In the standard implementation 
used for the internal needs of the router.