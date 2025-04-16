## Package auth
This package implements basic authentication logic. A more detailed authentication algorithm can be seen in [this test](https://github.com/uwine4850/foozy/blob/master/tests/auth/auth_test.go).

The authentication algorithm is as follows:
* Creating an authentication table using the __CreateAuthTable__ method.
* User registration using the __RegisterUser__ method.
* User login using the __LoginUser__ method.
* In order for the encryption keys to be updated, you need to use __globalflow__ and the function *bglobalflow.KeyUpdater(1)*.
    ```
    gf := globalflow.NewGlobalFlow(1000)
    gf.AddNotWaitTask(bglobalflow.KeyUpdater(1))
    gf.Run(manager)
    ```
* Now, when the keys are updated, you need to update the encoding with Middlewares. This can be done using a method *builtin_mddl.Auth*.
    ```
    mddl := middlewares.NewMiddleware()
    mddl.HandlerMddl(0, builtin_mddl.Auth("/login", mddlDb))
    ```

It is important to note that skipping the globalflow iteration without updating the coding, the user will not be able to update them and will need to log in again. This is because only the penultimate keys are matched, and to change the encoding you need to have matching keys. Accordingly, if the iteration is skipped, the necessary keys for encoding will simply be lost.

## Auth structure methods

__type AuthCookie struct__
The structure is intended to describe cookie data.

__type User struct__
The structure is intended to describe user data.

__RegisterUser__
```
RegisterUser(username string, password string) (int, error)
```
Creates a new user. If something goes wrong, it returns an error.
Returns the ID of the new user.

__LoginUser__
```
LoginUser(username string, password string) (*User, error)
```
Checks the correctness of user data. If there is no error, all the data is correct. 
Returns a __User__ structure that contains user data. Also adds cookie data about the user.

__AddAuthCookie__
```go
AddAuthCookie(uid string) error 
```
Adds the user's authentication cipher to the cookie.
It can then be used for authorization.

__UpdateAuthCookie__
```
UpdateAuthCookie(hashKey []byte, blockKey []byte, r *http.Request) error
```
Updates the encryption of user data stored in cookies. For a successful update, you need to enter the keys used for encryption in the *hashKey* and *blockKey* fields.

__ChangePassword__
```
ChangePassword(username string, oldPassword string, newPassword string) error
```
To change your password, you must enter your username, old password, and new password.

## Global functions for the auth package

__UserByUsername__
```
UserByUsername(db *database.Database, username string) (map[string]interface{}, error)
```
Checks if the user exists in the database. If it exists, returns its data.

__UserById__
```
UserById(db *database.Database, id any) (map[string]interface{}, error)
```
Searches for a user by ID. If it exists, returns its data.

__CreateAuthTable__
```
CreateAuthTable(database interfaces.IDatabase) error
```
Creates a table of users in the database. Used in the constructor of the ``Auth`` structure.

__HashPassword__
```
HashPassword(password string) (string, error)
```
Hashes the password for storage.

__ComparePassword__
```
ComparePassword(hashedPassword string, password string) error
```
Compares whether the hash of the password and the password match.

## Example of use
```go
dbArgs := database.DbArgs{
	Username: "root", Password: "1111", Host: "localhost", Port: "3408", DatabaseName: "foozy_test",
}
db := database.NewDatabase(dbArgs)
err := db.Connect()
if err != nil {
    panic(err)
}
defer func(db *database.Database) {
    err := db.Close()
    if err != nil {
	    panic(err)
	}
}(db)
_auth := auth.NewAuth(db)
err = _auth.RegisterUser("user", "111111")
if err != nil {
    panic(err)
}
user, err := _auth.LoginUser("user", "111111")
if err != nil {
    panic(err)
}
if err := _auth.AddAuthCookie(user.Id); err != nil {
    panic(err)
}
```