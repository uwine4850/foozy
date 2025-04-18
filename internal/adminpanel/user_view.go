package adminpanel

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbmapper"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/object"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
)

type UserViewObject struct {
	object.ObjView
}

func (uv *UserViewObject) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, uv.DB)
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

func (uv *UserViewObject) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	router.ServerError(w, err.Error(), manager)
}

func (uv *UserViewObject) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	rolesTableCreated, err := IsRolesTableCreated(uv.DB)
	if err != nil {
		return nil, err
	}
	if !rolesTableCreated {
		return nil, &ErrRolesTableNotCreated{}
	}
	roles, err := getAllRoles(uv.DB)
	if err != nil {
		return nil, err
	}
	id, _ := manager.OneTimeData().GetSlugParams(uv.Slug)
	userRole, err := UserRole(id, uv.DB)
	if err != nil {
		return nil, err
	}
	return object.Context{"roles": roles, "userRole": userRole}, nil
}

func UserView(db *database.Database) router.Handler {
	view := object.TemplateView{
		TemplatePath: filepath.Join(fpath.CurrentFileDir(), "templates/user_view.html"),
		View: &UserViewObject{
			object.ObjView{
				Name:       "user",
				DB:         db,
				TableName:  "auth",
				FillStruct: auth.User{},
				Slug:       "id",
			},
		},
	}
	return view.Call
}

type UserViewEditForm struct {
	UserId   []string `form:"user-id"`
	UserRole []string `form:"user-info-role-value"`
}

type UserViewEditFormObject struct {
	object.FormView
}

func (uv *UserViewEditFormObject) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, uv.DB)
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

func (uv *UserViewEditFormObject) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	router.ServerError(w, err.Error(), manager)
}

func (uv *UserViewEditFormObject) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.Context, error) {
	formInterface, err := uv.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	formObject := formInterface.(UserViewEditForm)
	if formObject.UserId[0] == "" {
		return nil, errors.New("error when editing user, ID not found")
	}
	if formObject.UserRole[0] == "---" || formObject.UserRole[0] == "" {
		_, err := qb.NewSyncQB(uv.DB.SyncQ()).Delete(USER_ROLES_TABLE).Where(qb.Compare("user_id", qb.EQUAL, formObject.UserId[0])).Exec()
		if err != nil {
			return nil, err
		}
	} else {
		role, err := qb.NewSyncQB(uv.DB.SyncQ()).SelectFrom("*", USER_ROLES_TABLE).Where(qb.Compare("user_id", qb.EQUAL, formObject.UserId[0])).Query()
		if err != nil {
			return nil, err
		}
		if len(role) > 0 && dbutils.ParseString(role[0]["role_name"]) != formObject.UserRole[0] {
			_, err := qb.NewSyncQB(uv.DB.SyncQ()).Update(USER_ROLES_TABLE, map[string]any{"role_name": formObject.UserRole[0]}).
				Where(qb.Compare("user_id", qb.EQUAL, formObject.UserId[0])).Exec()
			if err != nil {
				return nil, err
			}
		} else {
			_, err := qb.NewSyncQB(uv.DB.SyncQ()).
				Insert(USER_ROLES_TABLE, map[string]interface{}{"user_id": formObject.UserId[0], "role_name": formObject.UserRole[0]}).Exec()
			if err != nil {
				return nil, err
			}
		}
	}
	return object.Context{}, nil
}

func UserEditFormView(db *database.Database) router.Handler {
	view := object.TemplateRedirectView{
		RedirectUrl: "",
		View: &UserViewEditFormObject{
			object.FormView{
				FormStruct:       UserViewEditForm{},
				DB:               db,
				NotNilFormFields: []string{"user-info-role-value"},
				ValidateCSRF:     true,
			},
		},
	}
	return view.Call
}

func getAllRoles(db *database.Database) ([]RoleDB, error) {
	rolesDB, err := qb.NewSyncQB(db.SyncQ()).SelectFrom("*", ROLES_TABLE).Query()
	if err != nil {
		return nil, err
	}
	var roles []RoleDB
	mapper := dbmapper.NewMapper(rolesDB, typeopr.Ptr{}.New(&roles))
	if err := mapper.Fill(); err != nil {
		return nil, err
	}
	return roles, nil
}
