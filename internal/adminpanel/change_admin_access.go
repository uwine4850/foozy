package adminpanel

import (
	"errors"
	"net/http"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/object"
)

type AdminAccessEmptyForm struct{}

// Changes the way the admin page is accessed.
// Enables or disables access to the admin panel only with the administrator role.
// You can enable administrator-only access only if a user with this role exists.
type ChangeAdminAccessObject struct {
	object.FormView
}

func (aa *ChangeAdminAccessObject) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(aa.DB)
	if err != nil {
		return false, func() {
			router.ServerForbidden(w, manager)
		}
	}
	if ok {
		return true, func() {}
	} else {
		return false, func() {
			router.ServerForbidden(w, manager)
		}
	}
}

func (aa *ChangeAdminAccessObject) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	if errors.Is(err, &errAdminUserNotExist{}) {
		router.RedirectError(w, r, "/admin", err.Error())
		return
	}
	router.ServerError(w, err.Error(), manager)
}

func (aa *ChangeAdminAccessObject) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	admin_access_db, err := aa.DB.SyncQ().QB().Select("admin_access", ADMIN_SETTINGS_TABLE).Ex()
	if err != nil {
		return nil, err
	}
	admin_access_int, err := dbutils.ParseInt(admin_access_db[0]["admin_access"])
	if err != nil {
		return nil, err
	}
	var newAdminAccessValue int
	if admin_access_int == 0 {
		adminUserOK, err := AdminUserExists(aa.DB)
		if err != nil {
			return nil, err
		}
		if !adminUserOK {
			return object.ObjectContext{}, &errAdminUserNotExist{}
		}
		newAdminAccessValue = 1
	} else {
		newAdminAccessValue = 0
	}
	_, err = aa.DB.SyncQ().QB().Update(ADMIN_SETTINGS_TABLE, map[string]any{"admin_access": newAdminAccessValue}).Ex()
	if err != nil {
		return nil, err
	}
	return object.ObjectContext{}, nil
}

func ChangeAdminAccessView(db *database.Database) router.Handler {
	view := object.TemplateRedirectView{
		RedirectUrl: "/admin",
		View: &ChangeAdminAccessObject{
			FormView: object.FormView{
				FormStruct:       AdminAccessEmptyForm{},
				DB:               db,
				NotNilFormFields: []string{"*"},
				ValidateCSRF:     true,
			},
		},
	}
	return view.Call
}

type errAdminUserNotExist struct {
}

func (e *errAdminUserNotExist) Error() string {
	return "no user with the admin role"
}
