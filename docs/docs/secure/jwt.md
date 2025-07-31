## JWT
Working with JWT tokens.

#### NewHmacJwtWithClaims
Creating a JWT token. The library [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) is used for creation.
```golang
func NewHmacJwtWithClaims(claims jwt.Claims, manager interfaces.Manager) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.Key().Get32BytesKey().HashKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
```