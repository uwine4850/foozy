package adminpanel

import (
	"errors"
	"net/http"

	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/object"
)

type DeleteRoleForm struct {
	Name []string `form:"delete-name"`
}

// Role Deletion Form.
// Deletes a role by its ID. You cannot delete the administrator role.
type DeleteRoleObject struct {
	object.FormView
}

func (dr *DeleteRoleObject) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, dr.DB)
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

func (dr *DeleteRoleObject) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	if errors.Is(err, &errDeleteAdminRole{}) || IsObjectValidateCSRFError(err) {
		router.RedirectError(w, r, "/admin/users", err.Error())
	} else {
		router.ServerError(w, err.Error(), manager)
	}
}

func (dr *DeleteRoleObject) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	formObjectInterface, err := dr.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	formObject := formObjectInterface.(DeleteRoleForm)
	isAdmin, err := RoleIsAdmin(formObject.Name[0], dr.DB)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		return object.Context{}, &errDeleteAdminRole{}
	}
	_, err = qb.NewSyncQB(dr.DB.SyncQ()).Delete(ROLES_TABLE).Where(qb.Compare("name", qb.EQUAL, formObject.Name[0])).Exec()
	if err != nil {
		return nil, err
	}
	return object.Context{}, nil
}

func DeleteRoleView(db *database.Database) router.Handler {
	view := object.TemplateRedirectView{
		RedirectUrl: "/admin/users",
		View: &DeleteRoleObject{
			FormView: object.FormView{
				FormStruct:       DeleteRoleForm{},
				DB:               db,
				NotNilFormFields: []string{"*"},
				ValidateCSRF:     true,
			},
		},
	}
	return view.Call
}

type errDeleteAdminRole struct {
}

func (e *errDeleteAdminRole) Error() string {
	return "you can't delete the admin role"
}
