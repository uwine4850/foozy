package builtin_mddl

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

type onError func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)

// GenerateAndSetCsrf A middleware designed to generate a CSRF token. The token is set as a cookie value.
// To use it you need to run the method in a synchronous or asynchronous handler.
// maxAge - cookie lifetime.
// onError - a function that will be executed during an error.
func GenerateAndSetCsrf(maxAge int, onError onError) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, managerConfig interfaces.IManagerConfig) {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, managerConfig interfaces.IManagerConfig) {
		csrfCookie, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
		if err != nil || csrfCookie.Value == "" {
			csrfToken, err := GenerateCsrfToken()
			if err != nil {
				if onError != nil {
					onError(w, r, manager, err)
				}
				return
			}
			cookie := &http.Cookie{
				Name:     namelib.ROUTER.COOKIE_CSRF_TOKEN,
				Value:    csrfToken,
				MaxAge:   maxAge,
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
			}
			http.SetCookie(w, cookie)
			manager.Render().SetContext(map[string]interface{}{namelib.ROUTER.COOKIE_CSRF_TOKEN: fmt.Sprintf("<input name=\"%s\" type=\"hidden\" value=\"%s\">",
				namelib.ROUTER.COOKIE_CSRF_TOKEN, csrfToken)})
		}
	}
}

// GenerateCsrfToken generates a CSRF token.
func GenerateCsrfToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	csrfToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return csrfToken, nil
}
