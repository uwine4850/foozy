package cookies

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/secure"
)

func CreateSecureCookieData(hashKey []byte, blockKey []byte, w http.ResponseWriter, cookie *http.Cookie, cookieValue interface{}) error {
	secureValue, err := secure.CreateSecureData(hashKey, blockKey, cookieValue)
	if err != nil {
		return err
	}
	cookie.Value = secureValue
	http.SetCookie(w, cookie)
	return nil
}

func ReadSecureCookieData(hashKey []byte, blockKey []byte, r *http.Request, name string, readCookie interface{}) error {
	cookie, err := r.Cookie(name)
	if err != nil {
		return err
	}
	if err := secure.ReadSecureData(hashKey, blockKey, cookie.Value, &readCookie); err != nil {
		return err
	}
	return nil
}
