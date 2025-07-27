## Cookies
The package contains various interactions with standard cookies.

#### CreateSecureCookieData
CreateSecureCookieData creates cookie data using encoding and HMAC.<br>
hashKey is responsible for HMAC, and blockKey is for encoding.
```golang
func CreateSecureCookieData(hashKey []byte, blockKey []byte, w http.ResponseWriter, cookie *http.Cookie, cookieValue interface{}) error {
	secureValue, err := secure.CreateSecureData(hashKey, blockKey, cookieValue)
	if err != nil {
		return err
	}
	cookie.Value = secureValue
	http.SetCookie(w, cookie)
	return nil
}
```

#### ReadSecureCookieData
ReadSecureCookieData reads data encoded with [CreateSecureCookieData](#createsecurecookiedata).<br>
hashKey is responsible for HMAC, and blockKey is for encoding.
```golang
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
```

#### CreateSecureNoHMACCookieData
Creates an encoding of the cookie data, but without HMAC.
```golang
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
```

#### ReadSecureNoHMACCookieData
Reads data that was encoded using [CreateSecureNoHMACCookieData](#createsecurenohmaccookiedata).
```golang
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
```

#### SetStandartCookie
Just a handy feature for setting regular cookies.
```golang
func SetStandartCookie(w http.ResponseWriter, name string, value string, path string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}
```