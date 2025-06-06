package authtest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/builtin/bglobalflow"
	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
	"github.com/uwine4850/foozy/tests/common/tconf"
	testinitcnf "github.com/uwine4850/foozy/tests/common/test_init_cnf"
	"github.com/uwine4850/foozy/tests/common/tutils"
)

var mng = manager.NewManager(nil)
var newRouter = router.NewRouter(mng)

func TestMain(m *testing.M) {
	testinitcnf.InitCnf()
	_db := database.NewDatabase(tconf.DbArgs)
	if err := _db.Open(); err != nil {
		panic(err)
	}
	defer _db.Close()
	if err := database.InitDatabasePool(mng, _db); err != nil {
		panic(err)
	}

	if err := auth.CreateMysqlAuthTable(_db, tconf.DbArgs.DatabaseName); err != nil {
		panic(err)
	}
	_, err := _db.SyncQ().Query(fmt.Sprintf("TRUNCATE TABLE %s", namelib.AUTH.AUTH_TABLE))
	if err != nil {
		panic(err)
	}
	mng.Key().Generate32BytesKeys()

	mddlDb := database.NewDatabase(tconf.DbArgs)
	if err := mddlDb.Open(); err != nil {
		panic(err)
	}
	defer mddlDb.Close()
	mddl := middlewares.NewMiddleware()
	mddl.SyncMddl(0, builtin_mddl.Auth([]string{"/login"}, func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
		middlewares.SetMddlError(err, manager.OneTimeData())
	}))

	newRouter.Get("/register", register())
	newRouter.Get("/login", login())
	newRouter.Get("/uid", uid())
	newRouter.Get("/upd-keys", updKeys())
	newRouter.Get("/user-by-id", userById())
	serv := server.NewServer(tconf.PortAuth, newRouter, nil)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(tconf.PortAuth, 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
}

func register() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		dbRead, err := manager.Database().ConnectionPool(config.LoadedConfig().Default.Database.MainConnectionPoolName)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if err := auth.CreateMysqlAuthTable(dbRead, tconf.DbArgs.DatabaseName); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		au, err := auth.NewAuth(w, mng)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if _, err := au.RegisterUser("test", "111111"); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		return func() {}
	}
}

func login() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		au, err := auth.NewAuth(w, mng)
		if err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		if usr, err := au.LoginUser("test", "111111"); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		} else {
			err := au.AddAuthCookie(usr.Id)
			if err != nil {
				return func() { router.ServerError(w, err.Error(), manager) }
			}
		}
		return func() {}
	}
}

func uid() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		k := mng.Key().Get32BytesKey()
		var a auth.Cookie
		if err := cookies.ReadSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), r, namelib.AUTH.COOKIE_AUTH, &a); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		return func() {}
	}
}

func updKeys() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		k := mng.Key().Get32BytesKey()
		var a auth.Cookie
		if err := cookies.ReadSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), r, namelib.AUTH.COOKIE_AUTH, &a); err != nil {
			return func() { router.ServerError(w, err.Error(), manager) }
		}
		return func() {}
	}
}

func userById() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		dbRead, err := manager.Database().ConnectionPool(config.LoadedConfig().Default.Database.MainConnectionPoolName)
		if err != nil {
			panic(err)
		}
		user, err := auth.UserByID(dbRead, 1)
		if err != nil {
			panic(err)
		}
		if len(user) != 0 {
			id := user["id"]
			if id.(int64) == 1 {
				return func() { w.Write([]byte("OK")) }
			} else {
				return func() { w.Write([]byte("!OK")) }
			}
		}
		return func() { w.Write([]byte("!OK")) }
	}
}

func TestRegister(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortAuth, "register"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "" {
		t.Errorf("Error during registration: %s", string(body))
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestLogin(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortAuth, "login"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "" {
		t.Errorf("Error during sign in: %s", string(body))
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestReadUID(t *testing.T) {
	createReq, err := http.NewRequest("GET", tutils.MakeUrl(tconf.PortAuth, "login"), nil)
	if err != nil {
		t.Error(err)
	}
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Error(err)
	}
	defer createResp.Body.Close()

	readReq, err := http.NewRequest("GET", tutils.MakeUrl(tconf.PortAuth, "uid"), nil)
	if err != nil {
		t.Error(err)
	}

	readReq.AddCookie(createResp.Cookies()[0])

	readResp, err := http.DefaultClient.Do(readReq)
	if err != nil {
		t.Error(err)
	}
	defer readResp.Body.Close()

	body, err := io.ReadAll(readResp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "" {
		t.Errorf("Error during read UID: %s", string(body))
	}
}

func TestUpdKeys(t *testing.T) {
	k := mng.Key().Get32BytesKey()
	hashKey := k.HashKey()
	blockKey := k.BlockKey()
	gf := globalflow.NewGlobalFlow(1)
	gf.AddNotWaitTask(bglobalflow.KeyUpdater(1))
	gf.Run(mng)
	time.Sleep(2 * time.Second)
	if hashKey == k.HashKey() {
		t.Errorf("HashKey has not been updated.")
	}
	if blockKey == k.BlockKey() {
		t.Errorf("BlockKey has not been updated.")
	}

	get, err := http.Get(tutils.MakeUrl(tconf.PortAuth, "login"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "" {
		t.Errorf("Error during update key: %s", string(body))
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestUserByID(t *testing.T) {
	get, err := http.Get(tutils.MakeUrl(tconf.PortAuth, "user-by-id"))
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "OK" {
		t.Errorf("Error getting user by ID.")
	}
	err = get.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
