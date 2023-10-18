## Package auth
This package implements the basic authentication logic.

## Methods of the Auth structure
__RegisterUser__
```
RegisterUser(username string, password string) error
```
Creates a new user. If something goes wrong, it returns an error.

__LoginUser__
```
LoginUser(username string, password string) error
```
Checks the correctness of user data. If there is no error, all data is correct.

__ChangePassword__
```
ChangePassword(username string, oldPassword string, newPassword string) error
```
To change your password, you need to enter your username, old password, and new password.

__UserExist__
```
UserExist(username string) (map[string]interface{}, error)
```
Checks if the user is in the database. If it is, returns its data.

## Global packages for the auth package

__CreateAuthTable__
```
CreateAuthTable(database interfaces.IDatabase) error
```
Creates a table of users in the database. Used in the constructor of the ``Auth`` structure.
__HashPassword__
```
HashPassword(password string) (string, error)
```
Hash password for storage.

__ComparePassword__
```
ComparePassword(hashedPassword string, password string) error
```
Compares whether the password hash and password match.

## Example of use
```
db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
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
_auth, err := auth.NewAuth(db)
if err != nil {
    panic(err)
}
err = _auth.RegisterUser("user", "111111")
if err != nil {
    fmt.Println(err)
}
err = _auth.LoginUser("user", "111111")
if err != nil {
    fmt.Println(err)
	return
}
```