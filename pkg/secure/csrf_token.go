package secure

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
)

// ValidateFormCsrfToken checks the validity of the csrf token. If no errors are detected, the token is valid.
// It is desirable to use this method only after form.Parse() method.
func ValidateCookieCsrfToken(r *http.Request, token string) error {
	if token == "" {
		return ErrCsrfTokenNotFound{}
	}
	cookie, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
	if err != nil {
		return err
	}
	if cookie.Value != token {
		return ErrCsrfTokenDoesNotMatch{}
	}
	return nil
}

// ValidateHeaderCSRFToken validates the CSRF token based on its value in the header.
// For proper operation, the token must be set in cookies before verification.
func ValidateHeaderCSRFToken(r *http.Request, tokenName string) error {
	csrfToken := r.Header.Get(tokenName)
	if csrfToken == "" {
		return ErrCsrfTokenNotFound{}
	}
	return ValidateCookieCsrfToken(r, csrfToken)
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

func SetCSRFToken(maxAge int, httpOnly bool, w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
	csrfCookie, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
	if err != nil || csrfCookie.Value == "" {
		csrfToken, err := GenerateCsrfToken()
		if err != nil {
			return err
		}
		cookie := &http.Cookie{
			Name:     namelib.ROUTER.COOKIE_CSRF_TOKEN,
			Value:    csrfToken,
			MaxAge:   maxAge,
			HttpOnly: httpOnly,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		csrfHTMLString := fmt.Sprintf("<input name=\"%s\" type=\"hidden\" value=\"%s\">", namelib.ROUTER.COOKIE_CSRF_TOKEN, csrfToken)
		if manager.Render() != nil {
			manager.Render().SetContext(map[string]interface{}{namelib.ROUTER.COOKIE_CSRF_TOKEN: csrfHTMLString})
		}
		manager.OneTimeData().SetUserContext(namelib.ROUTER.COOKIE_CSRF_TOKEN, csrfHTMLString)
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
