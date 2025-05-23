package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Id int `json:"id"`
}

type Cookie struct {
	UID     int
	KeyDate time.Time
}

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
}

// Auth structure is designed to manage user authentication.
// It can be used to create a user, check the correctness of the login data, change the password and
// check the availability of the user.
type Auth struct {
	// database  *database.Database
	dbPool    interfaces.IReadDatabase
	tableName string
	w         http.ResponseWriter
	manager   interfaces.IManager
}

func NewAuth(w http.ResponseWriter, manager interfaces.IManager) (*Auth, error) {
	if !typeopr.IsPointer(manager) {
		panic("The manager config must be passed by pointer.")
	}
	db, err := manager.Database().ConnectionPool(config.LoadedConfig().Default.Database.MainConnectionPoolName)
	if err != nil {
		return nil, err
	}
	return &Auth{db, namelib.AUTH.AUTH_TABLE, w, manager}, nil
}

// RegisterUser registers the user in the database.
// It also checks the password and makes sure that there is no user with that login.
// Returns the ID of the new user.
func (a *Auth) RegisterUser(username string, password string) (int, error) {
	qUser := qb.NewSyncQB(a.dbPool.SyncQ()).SelectFrom("username", a.tableName).Where(
		qb.Compare("username", qb.EQUAL, username),
	).Limit(1)
	qUser.Merge()
	user, err := qUser.Query()
	if err != nil {
		return 0, err
	}
	if len(user) >= 1 {
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
	qIns := qb.NewSyncQB(a.dbPool.SyncQ()).Insert(a.tableName, map[string]interface{}{"username": username, "password": hashPass})
	qIns.Merge()
	insertCommand, err := qIns.Exec()
	if err != nil {
		return 0, err
	}
	insertUserId, ok := insertCommand["insertID"].(int64)
	if !ok {
		return 0, &ErrUserRegistration{}
	}
	return int(insertUserId), nil
}

// LoginUser check if the password and login are the same.
// If there was no error returns an [User] object with user data.
func (a *Auth) LoginUser(username string, password string) (*User, error) {
	userDB, err := UserByUsername(a.dbPool, username)
	if err != nil {
		return nil, err
	}
	if userDB == nil {
		return nil, ErrUserNotExist{username}
	}
	err = ComparePassword(dbutils.ParseString(userDB["password"]), password)
	if err != nil {
		return nil, err
	}
	var authItem User
	err = mapper.FillStructFromDb(&authItem, &userDB)
	if err != nil {
		return nil, err
	}
	return &authItem, nil
}

// Update Auth Cookie updates the cookie encoding.
// IMPORTANT: to work, you need to decode the data; accordingly, in the hashKey and blockKey fields you need to use the keys
// with which they were encoded. Next, the function itself will take new keys from ManagerConf.
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

// AddAuthCookie adds the user's authentication cipher to the cookie.
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

// ChangePassword changes the current user password.
func (a *Auth) ChangePassword(username string, oldPassword string, newPassword string) error {
	user, err := UserByUsername(a.dbPool, username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotExist{username}
	}
	if len(newPassword) < 6 {
		return ErrShortPassword{}
	}
	err = ComparePassword(dbutils.ParseString(user["password"]), oldPassword)
	if err != nil {
		return err
	}
	password, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	q := qb.NewSyncQB(a.dbPool.SyncQ()).Update(a.tableName, map[string]any{"password": password}).Where(
		qb.Compare("username", qb.EQUAL, username),
	)
	q.Merge()
	_, err = q.Exec()
	if err != nil {
		return err
	}
	return nil
}

// UserByUsername checks if the user is in the database.
// If it is found, it returns information about it.
func UserByUsername(db interfaces.IReadDatabase, username string) (map[string]interface{}, error) {
	qUser := qb.NewSyncQB(db.SyncQ()).SelectFrom("*", namelib.AUTH.AUTH_TABLE).Where(
		qb.Compare("username", qb.EQUAL, username),
	).Limit(1)
	qUser.Merge()
	user, err := qUser.Query()
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}
	return user[0], nil
}

// UserByID searches for a user by ID and returns it.
func UserByID(dbRead interfaces.IReadDatabase, id any) (map[string]interface{}, error) {
	qUser := qb.NewSyncQB(dbRead.SyncQ()).SelectFrom("*", namelib.AUTH.AUTH_TABLE).Where(
		qb.Compare("id", qb.EQUAL, id),
	).Limit(1)
	qUser.Merge()
	user, err := qUser.Query()
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}
	return user[0], nil
}

// CreateAuthTable creates a user authentication table.
func CreateAuthTable(dbRead interfaces.IReadDatabase, databaseName string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`.`%s` "+
		"(`id` INT NOT NULL AUTO_INCREMENT , "+
		"`username` VARCHAR(200) NOT NULL , "+
		"`password` TEXT NOT NULL , PRIMARY KEY (`id`))", databaseName, namelib.AUTH.AUTH_TABLE)
	_, err := dbRead.SyncQ().Query(sql)
	if err != nil {
		return err
	}
	return nil
}

// HashPassword generating a password hash from a string.
func HashPassword(password string) (string, error) {
	bytesPass := []byte(password)
	fromPassword, err := bcrypt.GenerateFromPassword(bytesPass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(fromPassword), nil
}

// ComparePassword check if the password hash and the password itself match.
func ComparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrPasswordsDontMatch{}
	}
	return nil
}
