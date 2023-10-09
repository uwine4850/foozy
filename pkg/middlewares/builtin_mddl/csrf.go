package builtin_mddl

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
)

func GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil || csrfCookie.Value == "" {
		csrfToken := utils.GenerateCsrfToken()
		cookie := &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			MaxAge:   1800,
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}
}
