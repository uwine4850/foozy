## package cookies
Cookie interactions. With this package you can set and read both 
encrypted and regular cookies.

__CreateSecureCookieData__
```go
CreateSecureCookieData(hashKey []byte, blockKey []byte, w http.ResponseWriter, cookie *http.Cookie, cookieValue interface{}) error
```
Creates secure cookies. Adds a signature using HMAC.

__ReadSecureCookieData__
```go
ReadSecureCookieData(hashKey []byte, blockKey []byte, r *http.Request, name string, readCookie interface{}) error
```
Reads data encoded with `CreateSecureCookieData`.

__CreateSecureNoHMACCookieData__
```go
CreateSecureNoHMACCookieData(key []byte, w http.ResponseWriter, cookie *http.Cookie, cookieValue interface{}) error
```
Creates an encoding of the cookie data, but without HMAC.

__ReadSecureNoHMACCookieData__
```go
ReadSecureNoHMACCookieData(key []byte, r *http.Request, name string, readValue interface{}) error
```
Reads data that was encoded using CreateSecureNoHMACCookieData.

__SetStandartCookie__
```go
SetStandartCookie(w http.ResponseWriter, name string, value string, path string, maxAge int)
```
Sets regular cookies without any non-standard encryption.