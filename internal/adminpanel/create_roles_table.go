package adminpanel

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
)

// Creates tables with roles.
// Three tables will be created if none exist:
//
//	Authentication table
//	Roles
//	User Role
//
// After the tables are created, the administrator role is created.
func CreateTables(db *database.Database) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		if err := db.Connect(); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		defer func() {
			if err := db.Close(); err != nil {
				router.ServerError(w, err.Error(), manager)
			}
		}()
		ok, err := AdminPermissions(db)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if !ok {
			return func() { router.ServerForbidden(w, manager) }
		}
		db.BeginTransaction()
		if err := createAuthTable(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}

		if err := CreateRolesTable(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}

		if err := CreateUserRolesTable(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}

		return func() { http.Redirect(w, r, "/admin/users", http.StatusFound) }
	}
}

func createAuthTable(db *database.Database) error {
	if err := auth.CreateAuthTable(db); err != nil {
		if err1 := db.RollBackTransaction(); err1 != nil {
			return err1
		}
		return err
	}
	return nil
}
