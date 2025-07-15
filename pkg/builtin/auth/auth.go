package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"golang.org/x/crypto/bcrypt"
)

// The AuthQuery interface is designed to abstract from the database during authentication.
// Any database can implement this interface and use it in the auth package.
type AuthQuery interface {
	// UserByUsername with username returns an [UnsafeUser] object that stores user data from the database.
	UserByUsername(username string) (*UnsafeUser, error)
	// UserById with id returns an [UnsafeUser] object that stores user data from the database.
	UserById(id any) (*UnsafeUser, error)
	// CreateNewUser creates a new user in the database.
	// The password must be encrypted in advance before passing it to this method.
	CreateNewUser(username string, hashPassword string) (result map[string]interface{}, err error)
	// ChangePassword changes the user's password.
	// The password must be encrypted in advance before passing it to this method.
	ChangePassword(userId string, newHashPassword string) (result map[string]interface{}, err error)
}

// MysqlAuthDatabase implementation of the [AuthQuery] interface.
// Executes database queries with specific methods only. Used only for the needs of the auth package.
//
// IMPORTANT: this structure uses the [UnsafeUser] object. Use this structure outside the auth package with caution.
type MysqlAuthDatabase struct {
	db        interfaces.DatabaseInteraction
	tableName string
}

func NewMysqlAuthQuery(db interfaces.DatabaseInteraction, tableName string) *MysqlAuthDatabase {
	return &MysqlAuthDatabase{
		db:        db,
		tableName: tableName,
	}
}

func (d *MysqlAuthDatabase) UserByUsername(username string) (*UnsafeUser, error) {
	qUser := qb.NewSyncQB(d.db.SyncQ()).SelectFrom("*", d.tableName).Where(
		qb.Compare("username", qb.EQUAL, username),
	).Limit(1)
	qUser.Merge()
	userQ, err := qUser.Query()
	if err != nil {
		return nil, err
	}
	var user UnsafeUser
	if len(userQ) == 0 {
		return nil, nil
	}
	if err := mapper.FillStructFromDb(&user, &userQ[0]); err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *MysqlAuthDatabase) UserById(id any) (*UnsafeUser, error) {
	qUser := qb.NewSyncQB(d.db.SyncQ()).SelectFrom("*", d.tableName).Where(
		qb.Compare("id", qb.EQUAL, id),
	).Limit(1)
	qUser.Merge()
	userQ, err := qUser.Query()
	if err != nil {
		return nil, err
	}
	var user UnsafeUser
	if len(userQ) == 0 {
		return nil, errors.New("user not exists")
	}
	if err := mapper.FillStructFromDb(&user, &userQ[0]); err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *MysqlAuthDatabase) CreateNewUser(username string, hashPassword string) (result map[string]interface{}, err error) {
	qIns := qb.NewSyncQB(d.db.SyncQ()).Insert(d.tableName, map[string]interface{}{"username": username, "password": hashPassword})
	qIns.Merge()
	result, err = qIns.Exec()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (d *MysqlAuthDatabase) ChangePassword(userId string, newHashPassword string) (result map[string]interface{}, err error) {
	q := qb.NewSyncQB(d.db.SyncQ()).Update(d.tableName, map[string]any{"password": newHashPassword}).Where(
		qb.Compare("id", qb.EQUAL, userId),
	)
	q.Merge()
	result, err = q.Exec()
	if err != nil {
		return nil, err
	}
	return result, nil
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Id int `json:"id"`
}

type Cookie struct {
	UID     int
	KeyDate time.Time
}

// UnsafeUser full user data.
// Not safe, should be used with caution.
type UnsafeUser struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

// User public user data that can be accessed by everyone.
type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
}

// Auth structure is designed to manage user authentication.
// It can be used to create a user, check the correctness of the login data, change the password and
// check the availability of the user.
type Auth struct {
	db        AuthQuery
	tableName string
	w         http.ResponseWriter
	manager   interfaces.Manager
}

func NewAuth(w http.ResponseWriter, db AuthQuery, manager interfaces.Manager) *Auth {
	return &Auth{db, namelib.AUTH.AUTH_TABLE, w, manager}
}

// RegisterUser registers the user in the database.
// It also checks the password and makes sure that there is no user with that login.
// Returns the ID of the new user.
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

// LoginUser check if the password and login are the same.
// If there was no error returns an [User] object with user data.
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

// UserByUsername checks if the user is in the database.
// If it is found, it returns information about it.
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

// UserByID searches for a user by ID and returns it.
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

// CreateMysqlAuthTable creates a user authentication table.
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
