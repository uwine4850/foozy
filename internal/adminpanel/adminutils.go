package adminpanel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbmapper"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/secure"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

const (
	ROLES_TABLE          = "roles"
	USER_ROLES_TABLE     = "user_roles"
	ADMIN_SETTINGS_TABLE = "admin_settings"
)

type RoleDB struct {
	Name string `db:"name"`
}

type UserRolesDB struct {
	Id       string `db:"id"`
	UserId   string `db:"user_id"`
	RoleName string `db:"role_name"`
}

// IsObjectValidateCSRFError сhecks if the CSRF validation of the token returns an error.
func IsObjectValidateCSRFError(err error) bool {
	if errors.Is(err, secure.ErrCsrfTokenNotFound{}) || errors.Is(err, secure.ErrCsrfTokenDoesNotMatch{}) {
		return true
	}
	return false
}

func AdminPermissions(db *database.Database) (bool, error) {
	isDebug := config.LoadedConfig().Default.Debug.Debug
	onlyAdminAccess, err := OnlyAdminAccess(db)
	if err != nil {
		return false, err
	}
	if onlyAdminAccess {
		adminExist, err := AdminUserExists(db)
		if err != nil {
			return false, err
		}
		if adminExist && isDebug {
			return true, nil
		}
	} else {
		if isDebug {
			return true, nil
		}
	}
	return false, nil
}

func AdminUserExists(db *database.Database) (bool, error) {
	roleId, err := db.SyncQ().QB().Select("role_name", USER_ROLES_TABLE).Where("role_name", "=", "ADMIN", "LIMIT 1").Ex()
	if err != nil {
		return false, err
	}
	if len(roleId) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func OnlyAdminAccess(db *database.Database) (bool, error) {
	isAdminSettingsTable, err := TableExists(db, ADMIN_SETTINGS_TABLE)
	if err != nil {
		return false, err
	}
	if isAdminSettingsTable {
		adminAccesDb, err := db.SyncQ().QB().Select("admin_access", ADMIN_SETTINGS_TABLE).Ex()
		if err != nil {
			return false, err
		}
		if len(adminAccesDb) > 0 {
			adminAccesInt, err := dbutils.ParseInt(adminAccesDb[0]["admin_access"])
			if err != nil {
				return false, err
			}
			if adminAccesInt == 1 {
				return true, nil
			}
		}
	}
	return false, nil
}

func TableExists(db *database.Database, tableName string) (bool, error) {
	res, err := db.SyncQ().Query(fmt.Sprintf("SHOW TABLES LIKE '%s';", tableName))
	if err != nil {
		return false, err
	}
	if len(res) > 0 {
		return true, nil
	}
	return false, nil
}

func UserRole(uid string, db *database.Database) (UserRolesDB, error) {
	userRolesDB, err := db.SyncQ().QB().Select("*", USER_ROLES_TABLE).Ex()
	if err != nil {
		return UserRolesDB{}, err
	}
	var userRoles []UserRolesDB
	mapper := dbmapper.NewMapper(userRolesDB, typeopr.Ptr{}.New(&userRoles))
	if err := mapper.Fill(); err != nil {
		return UserRolesDB{}, err
	}
	if len(userRoles) > 0 {
		return userRoles[0], nil
	} else {
		return UserRolesDB{}, nil
	}
}

func CreateSettingsTable(db *database.Database) error {
	_, err := db.SyncQ().Query(`
	CREATE TABLE IF NOT EXISTS admin_settings (
	id INT NOT NULL AUTO_INCREMENT,
	admin_access BOOLEAN NOT NULL DEFAULT FALSE,
	PRIMARY KEY (id));
	`)
	if err != nil {
		return err
	}
	res, err := db.SyncQ().QB().Select("id", ADMIN_SETTINGS_TABLE).Ex()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		_, err = db.SyncQ().Insert(ADMIN_SETTINGS_TABLE, map[string]any{})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateRolesTable(db *database.Database) error {
	_, err := db.SyncQ().Query(`
	CREATE TABLE IF NOT EXISTS roles (
	name VARCHAR(100) NOT NULL,
	PRIMARY KEY (name));
	`)
	if err != nil {
		if err1 := db.RollBackTransaction(); err1 != nil {
			return err1
		}
		return err
	}
	_, err = db.SyncQ().Insert(ROLES_TABLE, map[string]any{"name": "Admin"})
	if err != nil {
		if err1 := db.RollBackTransaction(); err1 != nil {
			return err1
		}
		return err
	}
	return nil
}

func CreateUserRolesTable(db *database.Database) error {
	_, err := db.SyncQ().Query(`
	CREATE TABLE IF NOT EXISTS user_roles (
	id INT NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	role_name VARCHAR(100) NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES auth(id) ON DELETE CASCADE,
	FOREIGN KEY (role_name) REFERENCES roles(name) ON DELETE CASCADE);
	`)
	if err != nil {
		if err1 := db.RollBackTransaction(); err1 != nil {
			return err1
		}
		return err
	}
	if err := db.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func RoleExists(name string, db *database.Database) (bool, error) {
	res, err := db.SyncQ().Query("SELECT * FROM roles WHERE UPPER(name) = UPPER(?);", name)
	if err != nil {
		return false, err
	}
	if len(res) > 0 {
		return true, nil
	}
	return false, nil
}

func RoleIsAdmin(name string, db *database.Database) (bool, error) {
	res, err := db.SyncQ().QB().Select("name", ROLES_TABLE).Where("name", "=", name).Ex()
	if err != nil {
		return false, err
	}
	if len(res) > 0 {
		if strings.EqualFold(dbutils.ParseString(res[0]["name"]), "Admin") {
			return true, nil
		}
	}
	return false, nil
}
