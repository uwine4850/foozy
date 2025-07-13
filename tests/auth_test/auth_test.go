package auth_test

import (
	"errors"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/builtin/bglobalflow"
	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/secure"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
	initcnf_t "github.com/uwine4850/foozy/tests1/common/init_cnf"
	"github.com/uwine4850/foozy/tests1/common/tutils"
	databasemock "github.com/uwine4850/foozy/tests1/database_mock"
)

type FakeAuthQuery struct{}

func (q *FakeAuthQuery) UserByUsername(username string) (*auth.UnsafeUser, error) {
	if username == "USER" {
		return &auth.UnsafeUser{
			Id:       1,
			Username: "USER",
			Password: "$2a$10$HBCPYsENH/sci4K9KJ7BXuJxSxguTtXCMfc.V2wQDq3ibCOT5QjG2",
		}, nil
	} else {
		return nil, nil
	}
}

func (q *FakeAuthQuery) UserById(id any) (*auth.UnsafeUser, error) {
	if id == "1" {
		return &auth.UnsafeUser{
			Id:       1,
			Username: "USER",
			Password: "$2a$10$HBCPYsENH/sci4K9KJ7BXuJxSxguTtXCMfc.V2wQDq3ibCOT5QjG2",
		}, nil
	} else {
		return nil, nil
	}
}

func (q *FakeAuthQuery) CreateNewUser(username string, hashPassword string) (result map[string]interface{}, err error) {
	return map[string]interface{}{"insertID": int64(2), "rowsAffected": 1}, err
}

func (q *FakeAuthQuery) ChangePassword(userId string, newHashPassword string) (result map[string]interface{}, err error) {
	return map[string]interface{}{"insertID": int64(1), "rowsAffected": 1}, err
}

var newManager = manager.NewManager(
	manager.NewOneTimeData(),
	nil,
	database.NewDatabasePool(),
)
var newMiddlewares = middlewares.NewMiddlewares()
var mockDatabase *databasemock.MysqlDatabase

func TestMain(m *testing.M) {
	initcnf_t.InitCnf()
	newManager.Key().Generate32BytesKeys()

	fakeAuthQuery := FakeAuthQuery{}
	syncQ := database.NewSyncQueries()
	asyncQ := database.NewAsyncQueries(syncQ)
	mockDatabase = databasemock.NewMysqlDatabase(syncQ, asyncQ)
	if err := mockDatabase.Open(); err != nil {
		panic(err)
	}

	excludePatterns := []string{"/test-register-user-exists", "/test-register", "/test-login", "/test-login-wrong-password",
		"/test-login-wrong-username", "/test-auth-jwt"}
	newMiddlewares.PreMiddleware(1,
		builtin_mddl.Auth(&fakeAuthQuery, excludePatterns,
			func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			},
		),
	)

	newAdapter := router.NewAdapter(newManager, newMiddlewares)
	newAdapter.SetOnErrorFunc(onError)
	newRouter := router.NewRouter(newAdapter)
	newRouter.Register(router.MethodGET, "/test-register-user-exists", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		myauth := auth.NewAuth(w, &fakeAuthQuery, manager)
		_, err := myauth.RegisterUser("USER", "111111")
		if err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-register", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		myauth := auth.NewAuth(w, &fakeAuthQuery, manager)
		uid, err := myauth.RegisterUser("USER_1", "111111")
		if err != nil {
			return err
		}
		if uid == 2 {
			w.Write([]byte("OK"))
		} else {
			w.Write([]byte("uid dont match"))
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-login", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		myauth := auth.NewAuth(w, &fakeAuthQuery, manager)
		user, err := myauth.LoginUser("USER", "111111")
		if err != nil {
			return err
		}
		if err := myauth.AddAuthCookie(user.Id); err != nil {
			return err
		}
		if user != nil && user.Id == 1 {
			w.Write([]byte("OK"))
		} else {
			w.Write([]byte("myauth.LoginUser returned an unexpected result"))
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-login-wrong-password", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		myauth := auth.NewAuth(w, &fakeAuthQuery, manager)
		_, err := myauth.LoginUser("USER", "WRONG PASSWORD")
		if err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-login-wrong-username", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		myauth := auth.NewAuth(w, &fakeAuthQuery, manager)
		_, err := myauth.LoginUser("WRONG USERNAME", "111111")
		if err != nil {
			return err
		}
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-auth-cookie", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		k := manager.Key().Get32BytesKey()
		var a auth.Cookie
		if err := cookies.ReadSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), r, namelib.AUTH.COOKIE_AUTH, &a); err != nil {
			return err
		}
		if a.UID == 0 {
			w.Write([]byte("auth cookie not found"))
		}
		w.Write([]byte("OK"))
		return nil
	})
	newRouter.Register(router.MethodGET, "/test-auth-jwt", func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		return nil
	})

	newServer := server.NewServer(tutils.PortAuth, newRouter, nil)
	go func() {
		if err := newServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	exitCode := m.Run()
	if err := newServer.Stop(); err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func onError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func TestRegisterUserExists(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-register-user-exists"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "User USER already exist." {
		if res != "" {
			t.Errorf("expected error not found: %s", res)
		} else {
			t.Error("expected error not found")
		}
	}
}

func TestRegister(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-register"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("register user error: %s", res)
	}
}

func TestLogin(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-login"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("login user error: %s", res)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-login-wrong-password"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "The passwords don't match." {
		if res != "" {
			t.Errorf("expected error not found: %s", res)
		} else {
			t.Error("expected error not found")
		}
	}
}

func TestLoginWrongUsername(t *testing.T) {
	resp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-login-wrong-username"))
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "User WRONG USERNAME not exist." {
		if res != "" {
			t.Errorf("expected error not found: %s", res)
		} else {
			t.Error("expected error not found")
		}
	}
}

func TestLoginAuthCookie(t *testing.T) {
	writeCookieResponse, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-login"))
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", tutils.MakeUrl(tutils.PortAuth, "test-auth-cookie"), nil)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(writeCookieResponse.Cookies()); i++ {
		req.AddCookie(writeCookieResponse.Cookies()[i])
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	res, err := tutils.ReadBody(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if res != "OK" {
		t.Errorf("expected error not found: %s", res)
	}
}

func TestUpdateKeys(t *testing.T) {
	k := newManager.Key().Get32BytesKey()
	hashKey := k.HashKey()
	blockKey := k.BlockKey()
	gf := globalflow.NewGlobalFlow(1)
	gf.AddNotWaitTask(bglobalflow.KeyUpdater(1))
	gf.Run(newManager)
	time.Sleep(2 * time.Second)
	if hashKey == k.HashKey() {
		t.Errorf("HashKey has not been updated.")
	}
	if blockKey == k.BlockKey() {
		t.Errorf("BlockKey has not been updated.")
	}
}

func TestAuthJWT(t *testing.T) {
	uidCalled := false
	updateCalled := false

	claims := auth.JWTClaims{
		Id: 1,
	}
	token, err := secure.NewHmacJwtWithClaims(claims, newManager)
	if err != nil {
		t.Error(err)
	}

	newMiddlewares.PreMiddleware(2, builtin_mddl.AuthJWT(
		func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (string, error) {
			return token, nil
		},
		func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, token string, AID int) error {
			updateCalled = true
			return nil
		},
		func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, AID int) error {
			uidCalled = true
			if AID != 1 {
				return errors.New("[currentID] error: uid does not match expectation")
			}
			return nil
		},
		func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error) {
			onError(w, r, err)
		},
	))

	uidResp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-auth-jwt"))
	if err != nil {
		t.Error(err)
	}
	uidRes, err := tutils.ReadBody(uidResp.Body)
	if err != nil {
		t.Error(err)
	}
	if uidRes != "" {
		t.Error(uidRes)
	}
	if !uidCalled {
		t.Error("[currentID] method is not called")
	}

	newManager.Key().Generate32BytesKeys()

	updateResp, err := http.Get(tutils.MakeUrl(tutils.PortAuth, "test-auth-jwt"))
	if err != nil {
		t.Error(err)
	}
	updateRes, err := tutils.ReadBody(updateResp.Body)
	if err != nil {
		t.Error(err)
	}
	if updateRes != "" {
		t.Error(updateRes)
	}

	if !updateCalled {
		t.Error("[updatedToken] method is not called")
	}
}

func TestAuthQCreateNewUser(t *testing.T) {
	authQ := auth.NewMysqlAuthQuery(mockDatabase, "auth")
	hashPassword, err := auth.HashPassword("111111")
	if err != nil {
		t.Error(err)
	}

	mockDatabase.Mock().
		ExpectExec(regexp.QuoteMeta("INSERT INTO auth ( username, password ) VALUES ( ?, ? )")).
		WithArgs("USER", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(10, 1))

	res, err := authQ.CreateNewUser("USER", hashPassword)
	if err != nil {
		t.Error(err)
	}
	if res["insertID"] != int64(10) && res["rowsAffected"] != int64(1) {
		t.Error("the result of the query does not match the expectation")
	}
}

func TestAuthQUserByUsername(t *testing.T) {
	authQ := auth.NewMysqlAuthQuery(mockDatabase, "auth")

	rows := sqlmock.NewRows([]string{"id", "username", "password"})
	rows.AddRow(1, "USER", "password")

	mockDatabase.Mock().
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM auth WHERE username = ?")).
		WithArgs("USER").
		WillReturnRows(rows)

	unsafeUser, err := authQ.UserByUsername("USER")
	if err != nil {
		t.Error(err)
	}
	expectUser := &auth.UnsafeUser{
		Id:       1,
		Username: "USER",
		Password: "password",
	}
	if !reflect.DeepEqual(unsafeUser, expectUser) {
		t.Error("object does not match expectations")
	}
}

func TestAuthQUserById(t *testing.T) {
	authQ := auth.NewMysqlAuthQuery(mockDatabase, "auth")

	rows := sqlmock.NewRows([]string{"id", "username", "password"})
	rows.AddRow(1, "USER", "password")

	mockDatabase.Mock().
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM auth WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(rows)

	unsafeUser, err := authQ.UserById(1)
	if err != nil {
		t.Error(err)
	}
	expectUser := &auth.UnsafeUser{
		Id:       1,
		Username: "USER",
		Password: "password",
	}
	if !reflect.DeepEqual(unsafeUser, expectUser) {
		t.Error("object does not match expectations")
	}
}

func TestAuthQChangePassword(t *testing.T) {
	authQ := auth.NewMysqlAuthQuery(mockDatabase, "auth")

	mockDatabase.Mock().
		ExpectExec(regexp.QuoteMeta("UPDATE auth SET password = ? WHERE id = ?")).
		WithArgs("new_password", "1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	res, err := authQ.ChangePassword("1", "new_password")
	if err != nil {
		t.Error(err)
	}
	if res["insertID"] != int64(0) && res["rowsAffected"] != int64(1) {
		t.Error("the result of the query does not match the expectation")
	}
}
