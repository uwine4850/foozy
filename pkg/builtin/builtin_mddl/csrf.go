package builtin_mddl

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/middlewares"
	"github.com/uwine4850/foozy/pkg/secure"
)

// GenerateAndSetCsrf A middleware designed to generate a CSRF token. The token is set as a cookie value.
// To use it you need to run the method in a synchronous or asynchronous handler.
// maxAge - cookie lifetime.
// onError - a function that will be executed during an error.
func GenerateAndSetCsrf(maxAge int, httpOnly bool) middlewares.PreMiddleware {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
		if err := secure.SetCSRFToken(maxAge, httpOnly, w, r, manager); err != nil {
			return err
		}
		return nil
	}
}
