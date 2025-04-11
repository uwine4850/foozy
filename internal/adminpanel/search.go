package adminpanel

import (
	"net/http"
	"path/filepath"

	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbmapper"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/object"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
)

type UserSearchForm struct {
	Id       []string `form:"id"`
	Username []string `form:"username"`
}

// Search for a user by ID.
type UserSearchByID struct {
	object.FormView
}

func (us *UserSearchByID) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, us.DB)
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

func (us *UserSearchByID) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	if IsObjectValidateCSRFError(err) {
		router.RedirectError(w, r, "/admin/users/search", err.Error())
	} else {
		router.ServerError(w, err.Error(), manager)
	}
}

func (us *UserSearchByID) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	formObjectInterface, err := us.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	formObject := formObjectInterface.(UserSearchForm)
	if len(formObject.Id) > 0 {
		res, err := qb.NewSyncQB(us.DB.SyncQ()).SelectFrom("*", namelib.AUTH.AUTH_TABLE).Where(qb.Compare("id", qb.EQUAL, formObject.Id[0])).Query()
		if err != nil {
			return nil, err
		}
		var out []auth.AuthUser
		mapper := dbmapper.NewMapper(res, typeopr.Ptr{}.New(&out))
		if err := mapper.Fill(); err != nil {
			return nil, err
		}
		manager.Render().SetContext(map[string]interface{}{"users": out, "search": "Search by ID: " + formObject.Id[0]})
	}
	return object.ObjectContext{}, nil
}

func UsersSearchByID(db *database.Database) router.Handler {
	view := object.TemplateView{
		TemplatePath: filepath.Join(fpath.CurrentFileDir(), "templates/search.html"),
		View: &UserSearchByID{
			FormView: object.FormView{
				FormStruct:    UserSearchForm{},
				DB:            db,
				NilIfNotExist: []string{"id", "username"},
				ValidateCSRF:  true,
			},
		},
	}
	return view.Call
}

// Search for a user by username.
type UserSearchByUsername struct {
	object.FormView
}

func (us *UserSearchByUsername) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	ok, err := AdminPermissions(r, manager, us.DB)
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

func (us *UserSearchByUsername) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	if IsObjectValidateCSRFError(err) {
		router.RedirectError(w, r, "/admin/users/search", err.Error())
	} else {
		router.ServerError(w, err.Error(), manager)
	}
}

func (us *UserSearchByUsername) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	formObjectInterface, err := us.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	formObject := formObjectInterface.(UserSearchForm)
	if len(formObject.Username) > 0 {
		res, err := qb.NewSyncQB(us.DB.SyncQ()).SelectFrom("*", namelib.AUTH.AUTH_TABLE).
			Where(qb.Compare("username", qb.LIKE, "'%"+formObject.Username[0]+"%'")).Limit(10).Query()
		if err != nil {
			return nil, err
		}
		var out []auth.AuthUser
		mapper := dbmapper.NewMapper(res, typeopr.Ptr{}.New(&out))
		if err := mapper.Fill(); err != nil {
			return nil, err
		}
		manager.Render().SetContext(map[string]interface{}{"users": out, "search": "Search by username: " + formObject.Username[0]})
	}
	return object.ObjectContext{}, nil
}

func UsersSearchByUsername(db *database.Database) router.Handler {
	view := object.TemplateView{
		TemplatePath: filepath.Join(fpath.CurrentFileDir(), "templates/search.html"),
		View: &UserSearchByUsername{
			FormView: object.FormView{
				FormStruct:    UserSearchForm{},
				DB:            db,
				NilIfNotExist: []string{"id", "username"},
				ValidateCSRF:  true,
			},
		},
	}
	return view.Call
}

func UsersSearchPage(db *database.Database) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		adminHTML := filepath.Join(fpath.CurrentFileDir(), "templates/search.html")
		manager.Render().SetTemplatePath(adminHTML)
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			router.ServerError(w, err.Error(), manager)
		}
		return func() {}
	}
}
