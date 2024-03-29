package auth

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"golang.org/x/crypto/bcrypt"
)

type ErrShortUsername struct {
}

func (receiver ErrShortUsername) Error() string {
	return "The username must be equal to or longer than 3 characters."
}

// Auth structure is designed to manage user authentication.
// It can be used to create a user, check the correctness of the login data, change the password and
// check the availability of the user.
type Auth struct {
	database  *database.Database
	tableName string
}

func NewAuth(database *database.Database) *Auth {
	return &Auth{database, "auth"}
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

// LoginUser User Login. This method only checks if the login details match.
func (a *Auth) LoginUser(username string, password string) error {
	user, err := a.UserExist(username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotExist{username}
	}
	err = ComparePassword(dbutils.ParseString(user["password"]), password)
	if err != nil {
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
	update, err := a.database.SyncQ().Update(a.tableName, []dbutils.DbEquals{{"password", password}}, dbutils.WHEquals(map[string]interface{}{"username": username}, "AND"))
	if err != nil {
		return err
	}
	fmt.Println(update)
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
