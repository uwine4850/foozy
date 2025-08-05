## auth middlewares

#### Auth
Auth is used to determine when to change the AUTH cookie encoding.<br>
When keys are changed, a change date is set. If the date does not match, then you need to change the encoding. 
It is important to note that only previous keys are saved; accordingly, it is impossible to update the encoding 
if two or more key iterations have passed, because the old keys are no longer known. 
This middleware should not work on the login page. Therefore, you need to specify the loginUrl correctly.

The `onErr` element is used for error management only within this middleware. When any error occurs, 
this function will be called instead of sending it to the router.
This is designed for more flexible control.

#### AuthJWT
Updates the JWT authentication encoding accordingly with key updates.<br> 
That is, the update depends directly on the frequency of key updates in GloablFlow.

The onErr element is used for error management only within this middleware. When any error occurs, 
this function will be called instead of sending it to the router. 
This is designed for more flexible control.

A more detailed explanation of the function arguments:

* Sets the JWT token for further work with it.
```golang
type SetToken func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (string, error)
```

* UpdatedToken function, which is called only if the token has been updated. Passes a single updated token.
```golang
type UpdatedToken func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, token string, AID int) error
```

* Function to which the user id is passed. This function is called each time the [AuthJWT] middleware is triggered. The function works both after token update and without update.
```golang
type CurrentUID func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, AID int) error
```