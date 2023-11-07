package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	ID    uint   `json:"id"`
}

type JwtC interface {
	CreateTokenUser(key string, claims *Claims, expiresAt time.Time) (string, error)
	ParseToken(key, token string, claims *Claims) (*jwt.Token, error)
}

type jwtC struct{}

func (j *jwtC) CreateTokenUser(key string, claims *Claims, expiresAt time.Time) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func (j *jwtC) ParseToken(key, token string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
}
