package builtin_mddl

import (
	"fmt"
	"net/http"
	"time"

	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
)

type OnError func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)

// Auth is used to determine when to change the AUTH cookie encoding.
// When keys are changed, a change date is set. If the date does not match, then you need to change the encoding.
// It is important to note that only previous keys are saved; accordingly, it is impossible to update the encoding
// if two or more key iterations have passed, because the old keys are no longer known.
// This middleware should not work on the login page. Therefore, you need to specify the loginUrl correctly.
func Auth(loginUrl string, db *database.Database, onErr OnError) middlewares.MddlFunc {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		pattern, ok := manager.OneTimeData().GetUserContext(namelib.ROUTER.URL_PATTERN)
		if !ok {
			onErr(w, r, manager, ErrUrlPatternNotExist{})
			return
		}
		if pattern == loginUrl {
			return
		}
		k := manager.Key().Get32BytesKey()
		var auth_date time.Time
		if err := cookies.ReadSecureNoHMACCookieData([]byte(k.StaticKey()), r, namelib.AUTH.COOKIE_AUTH_DATE, &auth_date); err != nil {
			onErr(w, r, manager, err)
			return
		}
		d1 := manager.Key().Get32BytesKey().Date().Format("02.01.2006 15:04:05")
		d2 := auth_date.Format("02.01.2006 15:04:05")
		if d1 != d2 {
			cc := database.NewConnectControl()
			if err := cc.OpenUnnamedConnection(db); err != nil {
				onErr(w, r, manager, err)
				return
			}
			defer func() {
				if err := cc.CloseAllUnnamedConnection(); err != nil {
					onErr(w, r, manager, err)
					return
				}
			}()
			_auth := auth.NewAuth(db, w, manager)
			if err := _auth.UpdateAuthCookie([]byte(k.OldHashKey()), []byte(k.OldBlockKey()), r); err != nil {
				onErr(w, r, manager, err)
				return
			}
			middlewares.SkipNextPageAndRedirect(manager.OneTimeData(), w, r, r.URL.Path)
		}
	}
}

type ErrUrlPatternNotExist struct {
}

func (e ErrUrlPatternNotExist) Error() string {
	return fmt.Sprintf("Data behind the %s key was not found.", namelib.ROUTER.URL_PATTERN)
}
