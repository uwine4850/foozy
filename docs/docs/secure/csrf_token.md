## CSRF token
The package specializes in working with CSRF tokens.

#### ValidateCookieCsrfToken
Checks the validity of the csrf token. If no errors are detected, the token is valid.
It is desirable to use this method only after [form.Parse](/router/form/form/#formparse) method.
```golang
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
```

#### ValidateHeaderCSRFToken
Validates the CSRF token based on its value in the header.
For proper operation, the token must be set in cookies before verification.
```golang
func ValidateHeaderCSRFToken(r *http.Request, tokenName string) error {
	csrfToken := r.Header.Get(tokenName)
	if csrfToken == "" {
		return ErrCsrfTokenNotFound{}
	}
	return ValidateCookieCsrfToken(r, csrfToken)
}
```

#### GenerateCsrfToken
Generates a CSRF token.
```golang
func GenerateCsrfToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	csrfToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return csrfToken, nil
}
```

#### SetCSRFToken
Sets the CSRF token in the cookie. Also, if [Render](/router/tmlengine/pagerender/) is previously installed in the [manager](/router/manager/manager), sets the template context `input` with the token by the key `namelib.ROUTER.COOKIE_CSRF_TOKEN`.
```golang
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
```