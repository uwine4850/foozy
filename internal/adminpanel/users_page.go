package adminpanel

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
)

// Renders the user and role interaction page.
func UsersPage(db *database.Database) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		if err := db.Connect(); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		defer func() {
			if err := db.Close(); err != nil {
				router.ServerError(w, err.Error(), manager)
			}
		}()
		ok, err := AdminPermissions(r, manager, db)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if !ok {
			return func() { router.ServerForbidden(w, manager) }
		}
		isRolesTableCreatedValue, err := isRolesTableCreated(db)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if isRolesTableCreatedValue {
			roles, err := getRoles(db)
			if err != nil {
				return func() { router.ServerError(w, err.Error(), manager) }
			}
			manager.Render().SetContext(map[string]interface{}{"roles": roles})
		}
		manager.Render().SetContext(map[string]interface{}{"rolesTableCreated": isRolesTableCreatedValue})
		router.CatchRedirectError(r, manager)
		adminHTML := filepath.Join(fpath.CurrentFileDir(), "templates/users.html")
		manager.Render().SetTemplatePath(adminHTML)
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		return func() {}
	}
}

func isRolesTableCreated(db *database.Database) (bool, error) {
	res, err := db.SyncQ().Query(fmt.Sprintf("SHOW TABLES LIKE '%s';", ROLES_TABLE))
	if err != nil {
		return false, err
	}
	if len(res) > 0 {
		return true, nil
	}
	return false, nil
}

func getRoles(db *database.Database) ([]RoleDB, error) {
	res, err := qb.NewSyncQB(db.SyncQ()).SelectFrom("*", ROLES_TABLE).Query()
	if err != nil {
		return nil, err
	}
	roles := make([]RoleDB, len(res))
	raw := mapper.NewDBRawStruct(&RoleDB{})
	if err := mapper.FillStructSliceFromDb(raw, &roles, &res); err != nil {
		return nil, err
	}
	return roles, nil
}
