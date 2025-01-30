package adminpanel

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
)

// CreateAdminSettingsTable creates a table with administration settings.
func CreateAdminSettingsTable(db *database.Database) router.Handler {
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
		if err := CreateSettingsTable(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		return func() { http.Redirect(w, r, "/admin", http.StatusFound) }
	}
}
