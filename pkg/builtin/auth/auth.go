package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type ErrShortUsername struct {
}

func (receiver ErrShortUsername) Error() string {
	return "The username must be equal to or longer than 3 characters."
}

type AuthCookie struct {
	UID     string
	KeyDate time.Time
}

type User struct {
	Id       string `db:"id"`
	Username string `db:"username"`
}

// Auth structure is designed to manage user authentication.
// It can be used to create a user, check the correctness of the login data, change the password and
// check the availability of the user.
type Auth struct {
	database  *database.Database
	tableName string
	w         http.ResponseWriter
	manager   interfaces.IManager
}

func NewAuth(database *database.Database, w http.ResponseWriter, manager interfaces.IManager) *Auth {
	if !utils.IsPointer(manager) {
		panic("The manager must be passed by pointer.")
	}
	return &Auth{database, "auth", w, manager}
}

// RegisterUser registers the user in the database.
// It also checks the password and makes sure that there is no user with that login.
func (a *Auth) RegisterUser(username string, password string) error {
	err := a.database.Ping()
	if err != nil {
		return err
	}
	user, err := a.database.SyncQ().Select([]string{"username"}, a.tableName, dbutils.WHEquals(map[string]interface{}{"username": username}, "AND"), 1)
	if err != nil {
		return err
	}
	if len(user) >= 1 {
		return ErrUserAlreadyExist{username}
	}
	if len(password) < 6 {
		return ErrShortPassword{}
	}
	if len(username) < 3 {
		return ErrShortUsername{}
	}
	hashPass, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = a.database.SyncQ().Insert(a.tableName, map[string]interface{}{"username": username, "password": hashPass})
	if err != nil {
		return err
	}
	return nil
}

// LoginUser check if the password and login are the same.
// Creates a cookie entry.
// Adds a USER variable to the user context, which contains user data from the auth table.
func (a *Auth) LoginUser(username string, password string) (*User, error) {
	userDB, err := a.UserExist(username)
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
	var user User
	if err := dbutils.FillStructFromDb(userDB, &user); err != nil {
		return nil, err
	}
	if err := a.addUserCookie(user.Id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *Auth) UpdateAuthCookie(hashKey []byte, blockKey []byte, r *http.Request) error {
	var authCookie AuthCookie
	if err := cookies.ReadSecureCookieData(hashKey, blockKey, r, "AUTH", &authCookie); err != nil {
		return err
	}
	if err := a.addUserCookie(authCookie.UID); err != nil {
		return err
	}
	return nil
}

func (a *Auth) addUserCookie(uid string) error {
	k := a.manager.Config().Get32BytesKey()
	if err := cookies.CreateSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), a.w, &http.Cookie{
		Name:     "AUTH",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}, &AuthCookie{UID: uid, KeyDate: a.manager.Get32BytesKey().Date()}); err != nil {
		return err
	}
	authDate := a.manager.Get32BytesKey().Date()
	if err := cookies.CreateSecureNoHMACCookieData([]byte(k.StaticKey()), a.w, &http.Cookie{
		Name:     "AUTH_DATE",
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
	user, err := a.UserExist(username)
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
	_, err = a.database.SyncQ().Update(a.tableName, []dbutils.DbEquals{{Name: "password", Value: password}}, dbutils.WHEquals(map[string]interface{}{"username": username}, "AND"))
	if err != nil {
		return err
	}
	return nil
}

// UserExist checks if the user is in the database.
func (a *Auth) UserExist(username string) (map[string]interface{}, error) {
	user, err := a.database.SyncQ().Select([]string{"*"}, a.tableName, dbutils.WHEquals(map[string]interface{}{"username": username}, "AND"), 1)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}
	return user[0], nil
}

// CreateAuthTable creates a user authentication table.
func CreateAuthTable(database *database.Database) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`.`auth` "+
		"(`id` INT NOT NULL AUTO_INCREMENT , "+
		"`username` VARCHAR(200) NOT NULL , "+
		"`password` TEXT NOT NULL , PRIMARY KEY (`id`))", database.DatabaseName())
	_, err := database.SyncQ().Query(sql)
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
