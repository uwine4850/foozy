## IManagerOneTimeData
This manager is responsible for transferring and saving the data of a separate request.
That is, it stores data only within one request, it cannot transmit 
data to other requests. Also, this manager can transfer data from Middleware 
to request

__SetUserContext__
```
SetUserContext(key string, value interface{})
```
Sets the custom context that is available within a valid entry.

__GetUserContext__
```
GetUserContext(key string) (any, bool)
```
Returns the user context. It is important to note that there may be 
messages from middleware, such as error values.

__DelUserContext__
```
DelUserContext(key string)
```
Deletes a custom context by key.

__SetSlugParams__
```
SetSlugParams(params map[string]string)
```
Sets the slug value of the parameter. In the standard implementation is used 
in the router.

__GetSlugParams__
```
GetSlugParams(key string) (string, bool)
```
Returns the slug value of the parameter.

__CreateNewManagerData__
```
CreateNewManagerData(manager interfaces.IManager) (interfaces.IManagerOneTimeData, error)
```
Creates and sets a new OneTimeData instance into the manager.