## auth
This package is designed for full interaction with the user account.

More detailed information on how to use this package can be found in the [tests](https://github.com/uwine4850/foozy/blob/master/tests/auth_test/auth_test.go).

### AuthQuery interface
Interface is designed to abstract from the database during authentication.<br>
Any database can implement this interface and use it in the auth package.

#### AuthQuery.UserByUsername
With username returns an [UnsafeUser](#unsafeuser) object that stores user data from the database.
```golang
UserByUsername(username string) (*UnsafeUser, error)
```

#### AuthQuery.UserById
With id returns an [UnsafeUser](#unsafeuser) object that stores user data from the database.
```golang
UserById(id any) (*UnsafeUser, error)
```

#### AuthQuery.CreateNewUser
Creates a new user in the database.<br>
The password must be encrypted in advance before passing it to this method.
```golang
CreateNewUser(username string, hashPassword string) (result map[string]interface{}, err error)
```

#### AuthQuery.ChangePassword
Changes the user's password.<br>
The password must be encrypted in advance before passing it to this method.
```golang
ChangePassword(userId string, newHashPassword string) (result map[string]interface{}, err error)
```

### UnsafeUser
Full user data.<br>
Not safe, should be used with caution.
```golang
type UnsafeUser struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
```

### User
Public user data that can be accessed by everyone.
```golang
type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
}
```

### Cookie
Presentation of authentication cookies.
```golang
type Cookie struct {
	UID     int
	KeyDate time.Time
}
```

### JWTClaims
Claims for JWT authentication.
```golang
type JWTClaims struct {
	jwt.RegisteredClaims
	Id int `json:"id"`
}
```

### Auth
Structure is designed to manage user authentication.<br>
It can be used to create a user, check the correctness of the login data, change the password and
check the availability of the user.

#### Auth.RegisterUser
Registers the user in the database.<br>
It also checks the password and makes sure that there is no user with that login.<br
Returns the ID of the new user.
```golang
func (a *Auth) RegisterUser(username string, password string) (userId int, err error) {
	user, err := a.db.UserByUsername(username)
	if err != nil {
		return 0, err
	}
	if user != nil {
		return 0, ErrUserAlreadyExist{username}
	}
	if len(password) < 6 {
		return 0, ErrShortPassword{}
	}
	if len(username) < 3 {
		return 0, ErrShortUsername{}
	}
	hashPass, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	result, err := a.db.CreateNewUser(username, hashPass)
	if err != nil {
		return 0, err
	}
	insertUserId, ok := result["insertID"].(int64)
	if !ok {
		return 0, &ErrUserRegistration{}
	}
	return int(insertUserId), nil
}
```

#### Auth.LoginUser
Check if the password and login are the same.<br>
If there was no error returns an [User](#user) object with user data.
```golang
func (a *Auth) LoginUser(username string, password string) (*User, error) {
	user, err := a.db.UserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExist{username}
	}
	err = ComparePassword(user.Password, password)
	if err != nil {
		return nil, err
	}
	return &User{
		user.Id,
		user.Username,
	}, nil
}
```

#### Auth.UpdateAuthCookie
Updates the cookie encoding.<br>
__IMPORTANT__: to work, you need to decode the data; accordingly, in the hashKey and blockKey fields you need to use the keys with which they were encoded.<br> Next, the function itself will take new keys from ManagerConf.
```golang
func (a *Auth) UpdateAuthCookie(hashKey []byte, blockKey []byte, r *http.Request) error {
	var authCookie Cookie
	if err := cookies.ReadSecureCookieData(hashKey, blockKey, r, namelib.AUTH.COOKIE_AUTH, &authCookie); err != nil {
		return err
	}
	if err := a.AddAuthCookie(authCookie.UID); err != nil {
		return err
	}
	return nil
}
```

#### Auth.AddAuthCookie
Adds the user's authentication cipher to the cookie.
```golang
func (a *Auth) AddAuthCookie(uid int) error {
	k := a.manager.Key().Get32BytesKey()
	if err := cookies.CreateSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), a.w, &http.Cookie{
		Name:     namelib.AUTH.COOKIE_AUTH,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}, &Cookie{UID: uid, KeyDate: a.manager.Key().Get32BytesKey().Date()}); err != nil {
		return err
	}
	authDate := a.manager.Key().Get32BytesKey().Date()
	if err := cookies.CreateSecureNoHMACCookieData([]byte(k.StaticKey()), a.w, &http.Cookie{
		Name:     namelib.AUTH.COOKIE_AUTH_DATE,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}, &authDate); err != nil {
		return err
	}
	return nil
}
```

#### Auth.ChangePassword
Changes the current user password.
```golang
func (a *Auth) ChangePassword(username string, oldPassword string, newPassword string) error {
	user, err := a.db.UserByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotExist{username}
	}
	if len(newPassword) < 6 {
		return ErrShortPassword{}
	}
	err = ComparePassword(user.Password, oldPassword)
	if err != nil {
		return err
	}
	password, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	if _, err := a.db.ChangePassword(username, password); err != nil {
		return err
	}
	return nil
}
```

#### UserByUsername
Checks if the user is in the database.<br>
If it is found, it returns information about it.
```golang
func UserByUsername(db AuthQuery, username string) (*User, error) {
	user, err := db.UserByUsername(username)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}
```

#### UserByID
Searches for a user by ID and returns it.
```golang
func UserByID(db AuthQuery, id any) (*User, error) {
	user, err := db.UserById(id)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}
```

#### CreateMysqlAuthTable
Creates a user authentication table.
```golang
func CreateMysqlAuthTable(dbInteraction interfaces.DatabaseInteraction, databaseName string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`.`%s` "+
		"(`id` INT NOT NULL AUTO_INCREMENT , "+
		"`username` VARCHAR(200) NOT NULL , "+
		"`password` TEXT NOT NULL , PRIMARY KEY (`id`))", databaseName, namelib.AUTH.AUTH_TABLE)
	_, err := dbInteraction.SyncQ().Query(sql)
	if err != nil {
		return err
	}
	return nil
}
```

#### HashPassword
Generating a password hash from a string.
```golang
func HashPassword(password string) (string, error) {
	bytesPass := []byte(password)
	fromPassword, err := bcrypt.GenerateFromPassword(bytesPass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(fromPassword), nil
}
```

#### ComparePassword
Check if the password hash and the password itself match.
```golang
func ComparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrPasswordsDontMatch{}
	}
	return nil
}
```