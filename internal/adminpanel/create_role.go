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

type CreateRoleForm struct {
	Name []string `form:"role-name"`
}

// User Role Creation Form.
// Handles POST method, you cannot create roles with the same name.
type CreateRoleObject struct {
	object.FormView
}

func (cr *CreateRoleObject) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, cr.DB)
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

func (cr *CreateRoleObject) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	if errors.Is(err, &errRoleAlreadyExists{}) || IsObjectValidateCSRFError(err) || errors.Is(err, &ErrRolesTableNotCreated{}) {
		router.RedirectError(w, r, "/admin/users", err.Error())
	} else {
		router.ServerError(w, err.Error(), manager)
	}
}

func (cr *CreateRoleObject) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	rolesTableCreated, err := IsRolesTableCreated(cr.DB)
	if err != nil {
		return nil, err
	}
	if !rolesTableCreated {
		return nil, &ErrRolesTableNotCreated{}
	}
	formObjectInterface, err := cr.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	formObject := formObjectInterface.(CreateRoleForm)
	roleFound, err := RoleExists(formObject.Name[0], cr.DB)
	if err != nil {
		return nil, err
	}
	if roleFound {
		return object.ObjectContext{}, &errRoleAlreadyExists{}
	}
	_, err = qb.NewSyncQB(cr.DB.SyncQ()).Insert(ROLES_TABLE, map[string]interface{}{"name": formObject.Name[0]}).Query()
	if err != nil {
		return nil, err
	}
	return object.ObjectContext{}, nil
}

func CreateRoleView(db *database.Database) router.Handler {
	view := object.TemplateRedirectView{
		RedirectUrl: "/admin/users",
		View: &CreateRoleObject{
			FormView: object.FormView{
				FormStruct:       CreateRoleForm{},
				DB:               db,
				NotNilFormFields: []string{"*"},
				ValidateCSRF:     true,
			},
		},
	}
	return view.Call
}

type errRoleAlreadyExists struct {
}

func (e *errRoleAlreadyExists) Error() string {
	return "the role already exists"
}
