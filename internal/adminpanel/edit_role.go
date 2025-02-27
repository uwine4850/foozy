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

type EditRoleForm struct {
	PrimaryName []string `form:"edit-role-primary-name"`
	Name        []string `form:"edit-role-name"`
}

// Role Edit Form.
// Edits the names of the role. You cannot edit the administrator role.
type EditRoleObject struct {
	object.FormView
}

func (er *EditRoleObject) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, er.DB)
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

func (er *EditRoleObject) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	if errors.Is(err, &errEditAdminRole{}) || IsObjectValidateCSRFError(err) {
		router.RedirectError(w, r, "/admin/users", err.Error())
	} else {
		router.ServerError(w, err.Error(), manager)
	}
}

func (er *EditRoleObject) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	formObjectInterface, err := er.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	formObject := formObjectInterface.(EditRoleForm)
	isAdmin, err := RoleIsAdmin(formObject.PrimaryName[0], er.DB)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		return object.ObjectContext{}, &errEditAdminRole{}
	}
	_, err = qb.NewSyncQB(er.DB.SyncQ()).Update(ROLES_TABLE, map[string]any{"name": formObject.Name[0]}).
		Where(qb.Compare("name", qb.EQUAL, formObject.PrimaryName[0])).Exec()
	if err != nil {
		return nil, err
	}
	return object.ObjectContext{}, nil
}

func EditRole(db *database.Database) router.Handler {
	view := object.TemplateRedirectView{
		RedirectUrl: "/admin/users",
		View: &EditRoleObject{
			FormView: object.FormView{
				FormStruct:       EditRoleForm{},
				DB:               db,
				NotNilFormFields: []string{"*"},
				ValidateCSRF:     true,
			},
		},
	}
	return view.Call
}

type errEditAdminRole struct {
}

func (e *errEditAdminRole) Error() string {
	return "you can't edit the admin role"
}
