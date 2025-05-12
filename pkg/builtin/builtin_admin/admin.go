package admin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/uwine4850/foozy/internal/adminpanel"
	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
)

var isStaticRegister = false

func AdminHandlerSet(db *database.Database, myRouter *router.Router) []map[string]map[string]router.Handler {
	adminSTATIC := filepath.Join(fpath.CurrentFileDir(), "../../../internal/adminpanel/static")
	if !isStaticRegister {
		myRouter.GetMux().Handle("/adminS/", http.StripPrefix("/adminS/", http.FileServer(http.Dir(adminSTATIC))))
		isStaticRegister = true
	}
	return []map[string]map[string]router.Handler{
		{
			router.GET: {"/admin/users": adminpanel.UsersPage(db)},
		},
		{
			router.GET: {"/admin/users/search": adminpanel.UsersSearchPage(db)},
		},
		{
			router.POST: {"/admin/users/search/search-id": adminpanel.UsersSearchByID(db)},
		},
		{
			router.POST: {"/admin/users/search/search-username": adminpanel.UsersSearchByUsername(db)},
		},
		{
			router.POST: {"/admin/users/create-roles-table": adminpanel.CreateTables(db)},
		},
		{
			router.POST: {"/admin/users/edit-role": adminpanel.EditRole(db)},
		},
		{
			router.POST: {"/admin/users/delete-role": adminpanel.DeleteRoleView(db)},
		},
		{
			router.POST: {"/admin/users/create-role": adminpanel.CreateRoleView(db)},
		},
		{
			router.POST: {"/admin/create-settings-table": adminpanel.CreateAdminSettingsTable(db)},
		},
		{
			router.POST: {"/admin/admin-access": adminpanel.ChangeAdminAccessView(db)},
		},
		{
			router.GET:  {"/admin/user/<id>": adminpanel.UserView(db)},
			router.POST: {"/admin/user/<id>": adminpanel.UserEditFormView(db)},
		},
	}
}

type AdminSettingsDB struct {
	Id          string `db:"id"`
	AdminAccess string `db:"admin_access"`
}

func AdminPage(db *database.Database) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		if err := db.Connect(); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		defer func() {
			if err := db.Close(); err != nil {
				router.ServerError(w, err.Error(), manager)
			}
		}()
		ok, err := adminpanel.AdminPermissions(r, manager, db)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if !ok {
			return func() { router.ServerForbidden(w, manager) }
		}
		adminHTML := filepath.Join(fpath.CurrentFileDir(), "../../../internal/adminpanel/templates/admin.html")
		settingsTableCreated, err := IsSettingsTableCreated(db)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		manager.Render().SetContext(map[string]interface{}{"isSettingsTableCreated": settingsTableCreated})

		if settingsTableCreated {
			adminSettingsDB, err := AdminSettings(db)
			if err != nil {
				fmt.Println(reflect.TypeOf(err))
				return func() { router.ServerError(w, err.Error(), manager) }
			}
			manager.Render().SetContext(map[string]interface{}{"adminSettingsDB": adminSettingsDB})
		}
		router.CatchRedirectError(r, manager)
		manager.Render().SetTemplatePath(adminHTML)
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			router.ServerError(w, err.Error(), manager)
		}
		return func() {}
	}
}

func IsSettingsTableCreated(db *database.Database) (bool, error) {
	res, err := db.SyncQ().Query(fmt.Sprintf("SHOW TABLES LIKE '%s';", adminpanel.ADMIN_SETTINGS_TABLE))
	if err != nil {
		return false, err
	}
	if len(res) > 0 {
		return true, nil
	}
	return false, nil
}

func AdminSettings(db *database.Database) (AdminSettingsDB, error) {
	res, err := qb.NewSyncQB(db.SyncQ()).SelectFrom("*", adminpanel.ADMIN_SETTINGS_TABLE).Query()
	if err != nil {
		return AdminSettingsDB{}, err
	}
	if len(res) == 0 {
		return AdminSettingsDB{}, nil
	}
	adminSettingsDB := make([]AdminSettingsDB, len(res))
	if err := mapper.FillStructSliceFromDb(&adminSettingsDB, &res); err != nil {
		return AdminSettingsDB{}, nil
	}
	return adminSettingsDB[0], nil
}

func IsAuthUserAdmin(r *http.Request, mng interfaces.IManager, db *database.Database) (bool, error) {
	return adminpanel.UserIsAdmin(r, mng, db)
}

func CheckRole(userID string, roleName string, db *database.Database) (bool, error) {
	ex, err := qb.SelectExists(qb.NewSyncQB(db.SyncQ()), adminpanel.USER_ROLES_TABLE,
		qb.Compare("user_id", qb.EQUAL, userID), qb.AND,
		qb.Compare("role_name", qb.EQUAL, roleName))
	if err != nil {
		return false, err
	}
	return ex, nil
}
