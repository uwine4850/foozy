package builtin_mddl

import (
	"fmt"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/utils"
)

// GenerateAndSetCsrf A middleware designed to generate a CSRF token. The token is set as a cookie value.
// To use it you need to run the method in a synchronous or asynchronous handler.
func GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	csrfCookie, err := r.Cookie(namelib.CSRF_TOKEN_COOKIE)
	if err != nil || csrfCookie.Value == "" {
		csrfToken, err := utils.GenerateCsrfToken()
		if err != nil {
			router.ServerError(w, err.Error(), manager.Config())
			return
		}
		cookie := &http.Cookie{
			Name:     namelib.CSRF_TOKEN_COOKIE,
			Value:    csrfToken,
			MaxAge:   1800,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		manager.Render().SetContext(map[string]interface{}{namelib.CSRF_TOKEN_COOKIE: fmt.Sprintf("<input name=\"%s\" type=\"hidden\" value=\"%s\">",
			namelib.CSRF_TOKEN_COOKIE, csrfToken)})
	}
}
