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

// Auth is used to determine when to change the AUTH cookie encoding.
// When keys are changed, a change date is set. If the date does not match, then you need to change the encoding.
// It is important to note that only previous keys are saved; accordingly, it is impossible to update the encoding
// if two or more key iterations have passed, because the old keys are no longer known.
func Auth(loginUrl string, db *database.Database) middlewares.MddlFunc {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		pattern, ok := manager.OneTimeData().GetUserContext(namelib.URL_PATTERN)
		if !ok {
			middlewares.SetMddlError(ErrUrlPatternNotExist{}, manager.OneTimeData())
			return
		}
		if pattern == loginUrl {
			return
		}
		k := manager.Config().Get32BytesKey()
		var auth_date time.Time
		if err := cookies.ReadSecureNoHMACCookieData([]byte(k.StaticKey()), r, namelib.AUTH_DATE_COOKIE, &auth_date); err != nil {
			middlewares.SetMddlError(err, manager.OneTimeData())
			return
		}
		d1 := manager.Config().Get32BytesKey().Date().Format("02.01.2006 15:04:05")
		d2 := auth_date.Format("02.01.2006 15:04:05")
		if d1 != d2 {
			cc := database.NewConnectControl()
			if err := cc.OpenUnnamedConnection(db); err != nil {
				middlewares.SetMddlError(err, manager.OneTimeData())
				return
			}
			defer func() {
				if err := cc.CloseAllUnnamedConnection(); err != nil {
					middlewares.SetMddlError(err, manager.OneTimeData())
					return
				}
			}()
			_auth := auth.NewAuth(db, w, manager)
			if err := _auth.UpdateAuthCookie([]byte(k.OldHashKey()), []byte(k.OldBlockKey()), r); err != nil {
				middlewares.SetMddlError(err, manager.OneTimeData())
				return
			}
			middlewares.SkipNextPageAndRedirect(manager.OneTimeData(), w, r, r.URL.Path)
		}
	}
}

type ErrUrlPatternNotExist struct {
}

func (e ErrUrlPatternNotExist) Error() string {
	return fmt.Sprintf("Data behind the %s key was not found.", namelib.URL_PATTERN)
}
