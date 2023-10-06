package builtin_mddl

import (
	"github.com/uwine4850/foozy/internal/interfaces"
	"github.com/uwine4850/foozy/internal/utils"
	"net/http"
)

func GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	csrfToken := utils.GenerateCsrfToken()
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		MaxAge:   1800,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
