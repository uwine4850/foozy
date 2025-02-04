package builtin_mddl

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/secure"
)

type onError func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)

// GenerateAndSetCsrf A middleware designed to generate a CSRF token. The token is set as a cookie value.
// To use it you need to run the method in a synchronous or asynchronous handler.
// maxAge - cookie lifetime.
// onError - a function that will be executed during an error.
func GenerateAndSetCsrf(maxAge int, onError onError) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	return func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		if err := secure.SetCSRFToken(maxAge, w, r, manager); err != nil {
			onError(w, r, manager, err)
		}
	}
}
