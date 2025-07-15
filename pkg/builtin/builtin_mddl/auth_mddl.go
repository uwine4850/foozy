package builtin_mddl

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/cookies"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/secure"
)

type OnError func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error)

// Auth is used to determine when to change the AUTH cookie encoding.
// When keys are changed, a change date is set. If the date does not match, then you need to change the encoding.
// It is important to note that only previous keys are saved; accordingly, it is impossible to update the encoding
// if two or more key iterations have passed, because the old keys are no longer known.
// This middleware should not work on the login page. Therefore, you need to specify the loginUrl correctly.
//
// The onErr element is used for error management only within this middleware. When any error occurs,
// this function will be called instead of sending it to the router.
// This is designed for more flexible control.
func Auth(adb auth.AuthQuery, excludePatterns []string, onErr OnError) middlewares.PreMiddleware {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		pattern, ok := manager.OneTimeData().GetUserContext(namelib.ROUTER.URL_PATTERN)
		if !ok {
			onErr(w, r, manager, ErrUrlPatternNotExist{})
			return middlewares.ErrStopMiddlewares{}
		}
		if slices.Contains(excludePatterns, pattern.(string)) {
			return nil
		}
		k := manager.Key().Get32BytesKey()
		var auth_date time.Time
		if err := cookies.ReadSecureNoHMACCookieData([]byte(k.StaticKey()), r, namelib.AUTH.COOKIE_AUTH_DATE, &auth_date); err != nil {
			onErr(w, r, manager, err)
			return middlewares.ErrStopMiddlewares{}
		}
		d1 := manager.Key().Get32BytesKey().Date().Format("02.01.2006 15:04:05")
		d2 := auth_date.Format("02.01.2006 15:04:05")
		if d1 != d2 {
			_auth := auth.NewAuth(w, adb, manager)
			if err := _auth.UpdateAuthCookie([]byte(k.OldHashKey()), []byte(k.OldBlockKey()), r); err != nil {
				onErr(w, r, manager, err)
				return middlewares.ErrStopMiddlewares{}
			}
			middlewares.SkipNextPageAndRedirect(manager.OneTimeData(), w, r, r.URL.Path)
		}
		return nil
	}
}

// SetToken sets the JWT token for further work with it.
type SetToken func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (string, error)

// UpdatedToken function, which is called only if the token has been updated.
// Passes a single updated token.
type UpdatedToken func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, token string, AID int) error

// CurrentUID function to which the user id is passed.
// This function is called each time the [AuthJWT] middleware is triggered.
// The function works both after token update and without update.
type CurrentUID func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, AID int) error

// AuthJWT updates the JWT authentication encoding accordingly with key updates.
// That is, the update depends directly on the frequency of key updates in GloablFlow.
//
// The onErr element is used for error management only within this middleware. When any error occurs,
// this function will be called instead of sending it to the router.
// This is designed for more flexible control.
func AuthJWT(setToken SetToken, updatedToken UpdatedToken, currentUID CurrentUID, onErr OnError) middlewares.PreMiddleware {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		tokenString, err := setToken(w, r, manager)
		if err != nil {
			onErr(w, r, manager, err)
			return nil
		}
		if tokenString == "" {
			return nil
		}
		_claims := &auth.JWTClaims{}
		_, err = jwt.ParseWithClaims(tokenString, _claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(manager.Key().HashKey()), nil
		})

		// The token cannot be decrypted by the current key.
		// This means that it has been changed and you should try the previous key.
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			newClaims := &auth.JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, newClaims, func(t *jwt.Token) (interface{}, error) {
				return []byte(manager.Key().OldHashKey()), nil
			})
			if err != nil {
				onErr(w, r, manager, err)
				return nil
			}
			// If the previous key fits and the token is valid, you need to update the encoding.
			if token.Valid {
				updatedTokenString, err := secure.NewHmacJwtWithClaims(newClaims, manager)
				if err != nil {
					onErr(w, r, manager, err)
					return nil
				}
				if err := updatedToken(w, r, manager, updatedTokenString, newClaims.Id); err != nil {
					onErr(w, r, manager, err)
					return nil
				}
				if err := currentUID(w, r, manager, newClaims.Id); err != nil {
					onErr(w, r, manager, err)
					return nil
				}
				return nil
			} else {
				onErr(w, r, manager, &ErrJWTNotValid{})
				return nil
			}
		}
		if err := currentUID(w, r, manager, _claims.Id); err != nil {
			onErr(w, r, manager, err)
			return nil
		}
		return nil
	}
}

type ErrUrlPatternNotExist struct {
}

func (e ErrUrlPatternNotExist) Error() string {
	return fmt.Sprintf("Data behind the %s key was not found.", namelib.ROUTER.URL_PATTERN)
}

type ErrJWTNotValid struct{}

func (e ErrJWTNotValid) Error() string {
	return "JWT is not valid"
}
