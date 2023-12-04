package builtin_mddl

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
)

// GenerateAndSetCsrf A middleware designed to generate a CSRF token. The token is set as a cookie value.
// To use it you need to run the method in a synchronous or asynchronous handler.
func GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData) {
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil || csrfCookie.Value == "" {
		csrfToken := utils.GenerateCsrfToken()
		cookie := &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			MaxAge:   1800,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}
}
