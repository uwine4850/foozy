package cookies

import (
	"encoding/json"
	"net/http"

	"github.com/uwine4850/foozy/pkg/secure"
)

// CreateSecureCookieData creates cookie data using encoding and HMAC.
// hashKey is responsible for HMAC, and blockKey is for encoding.
func CreateSecureCookieData(hashKey []byte, blockKey []byte, w http.ResponseWriter, cookie *http.Cookie, cookieValue interface{}) error {
	secureValue, err := secure.CreateSecureData(hashKey, blockKey, cookieValue)
	if err != nil {
		return err
	}
	cookie.Value = secureValue
	http.SetCookie(w, cookie)
	return nil
}

// ReadSecureCookieData reads data encoded with CreateSecureCookieData.
// hashKey is responsible for HMAC, and blockKey is for encoding.
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

// CreateSecureNoHMACCookieData creates an encoding of the cookie data, but without HMAC.
func CreateSecureNoHMACCookieData(key []byte, w http.ResponseWriter, cookie *http.Cookie, cookieValue interface{}) error {
	data, err := json.Marshal(cookieValue)
	if err != nil {
		return err
	}
	enc, err := secure.Encrypt(key, data)
	if err != nil {
		return err
	}
	cookie.Value = enc
	http.SetCookie(w, cookie)
	return nil
}

// ReadSecureNoHMACCookieData reads data that was encoded using CreateSecureNoHMACCookieData.
func ReadSecureNoHMACCookieData(key []byte, r *http.Request, name string, readValue interface{}) error {
	cookie, err := r.Cookie(name)
	if err != nil {
		return err
	}
	dec, err := secure.Decrypt(key, cookie.Value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dec, readValue); err != nil {
		return err
	}
	return nil
}
