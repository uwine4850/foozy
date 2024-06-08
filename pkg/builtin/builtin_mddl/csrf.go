package builtin_mddl

import (
	"fmt"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/utils"
)

// GenerateAndSetCsrf A middleware designed to generate a CSRF token. The token is set as a cookie value.
// To use it you need to run the method in a synchronous or asynchronous handler.
func GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil || csrfCookie.Value == "" {
		csrfToken, err := utils.GenerateCsrfToken()
		if err != nil {
			router.ServerError(w, err.Error(), manager)
			return
		}
		cookie := &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			MaxAge:   1800,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		manager.SetContext(map[string]interface{}{"csrf_token": fmt.Sprintf("<input name=\"csrf_token\" type=\"hidden\" value=\"%s\">", csrfToken)})
	}
}
