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
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/router/tmlengine"
	"github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
)

var dbArgs = database.DbArgs{
	Username: "root", Password: "1111", Host: "localhost", Port: "3408", DatabaseName: "foozy_test",
}

var mng = manager.NewManager(nil)
var newRouter = router.NewRouter(mng)

func TestMain(m *testing.M) {
	_db := database.NewDatabase(dbArgs)
	if err := _db.Connect(); err != nil {
		panic(err)
	}
	_, err := _db.SyncQ().Query(fmt.Sprintf("DELETE FROM %s", namelib.AUTH_TABLE))
	if err != nil {
		panic(err)
	}
	_db.Close()
	mng.Config().Generate32BytesKeys()
	mng.Config().Debug(true)

	mddlDb := database.NewDatabase(dbArgs)
	if err := mddlDb.Connect(); err != nil {
		panic(err)
	}
	defer mddlDb.Close()
	mddl := middlewares.NewMiddleware()
	mddl.HandlerMddl(0, builtin_mddl.Auth("/login", mddlDb))

	newRouter.SetTemplateEngine(&tmlengine.TemplateEngine{})
	newRouter.Get("/register", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		db := database.NewDatabase(dbArgs)
		cc := database.NewConnectControl()
		if err := cc.OpenUnnamedConnection(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		defer func() {
			if err := cc.CloseAllUnnamedConnection(); err != nil {
				router.ServerError(w, err.Error(), manager.Config())
			}
		}()
		if err := auth.CreateAuthTable(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		au := auth.NewAuth(db, w, manager)
		if err := au.RegisterUser("test", "111111"); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		return func() {}
	})
	newRouter.Get("/login", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		db := database.NewDatabase(dbArgs)
		cc := database.NewConnectControl()
		if err := cc.OpenUnnamedConnection(db); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		defer func() {
			if err := cc.CloseAllUnnamedConnection(); err != nil {
				router.ServerError(w, err.Error(), manager.Config())
			}
		}()
		au := auth.NewAuth(db, w, manager)
		if _, err := au.LoginUser("test", "111111"); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		return func() {}
	})
	newRouter.Get("/uid", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		k := manager.Config().Get32BytesKey()
		var a auth.AuthCookie
		if err := cookies.ReadSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), r, namelib.AUTH_COOKIE, &a); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		return func() {}
	})
	newRouter.Get("/upd-keys", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		k := manager.Config().Get32BytesKey()
		var a auth.AuthCookie
		if err := cookies.ReadSecureCookieData([]byte(k.HashKey()), []byte(k.BlockKey()), r, namelib.AUTH_COOKIE, &a); err != nil {
			return func() { router.ServerError(w, err.Error(), manager.Config()) }
		}
		return func() {}
	})
	serv := server.NewServer(":8060", newRouter)
	go func() {
		err = serv.Start()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
	if err := server.WaitStartServer(":8060", 5); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	err = serv.Stop()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestRegister(t *testing.T) {
	get, err := http.Get("http://localhost:8060/register")
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
	get, err := http.Get("http://localhost:8060/login")
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
	createReq, err := http.NewRequest("GET", "http://localhost:8060/login", nil)
	if err != nil {
		t.Error(err)
	}
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Error(err)
	}
	defer createResp.Body.Close()

	readReq, err := http.NewRequest("GET", "http://localhost:8060/uid", nil)
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
	k := mng.Config().Get32BytesKey()
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

	get, err := http.Get("http://localhost:8060/login")
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
