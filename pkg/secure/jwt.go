package secure

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

func NewHmacJwtWithClaims(claims jwt.Claims, manager interfaces.Manager) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.Key().Get32BytesKey().HashKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
