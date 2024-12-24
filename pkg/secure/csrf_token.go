package secure

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

// ValidateHeaderCSRFToken validates the CSRF token based on its value in the header.
// For proper operation, the token must be set in cookies before verification.
func ValidateHeaderCSRFToken(r *http.Request, tokenName string) error {
	csrfToken := r.Header.Get(tokenName)
	if csrfToken == "" {
		return ErrCsrfTokenNotFound{}
	}
	cookie, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
	if err != nil {
		return err
	}
	if cookie.Value != csrfToken {
		return ErrCsrfTokenDoesNotMatch{}
	}
	return nil
}

// ValidateFormCsrfToken checks the validity of the csrf token. If no errors are detected, the token is valid.
// It is desirable to use this method only after form.Parse() method.
func ValidateFormCsrfToken(r *http.Request, frm interfaces.IForm) error {
	csrfToken := frm.Value(namelib.ROUTER.COOKIE_CSRF_TOKEN)
	if csrfToken == "" {
		return ErrCsrfTokenNotFound{}
	}
	cookie, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
	if err != nil {
		return err
	}
	if cookie.Value != csrfToken {
		return ErrCsrfTokenDoesNotMatch{}
	}
	return nil
}

type ErrCsrfTokenNotFound struct {
}

func (receiver ErrCsrfTokenNotFound) Error() string {
	return "Csrf token not found."
}

type ErrCsrfTokenDoesNotMatch struct {
}

func (receiver ErrCsrfTokenDoesNotMatch) Error() string {
	return "Csrf token does not match. The validity time may have expired."
}
