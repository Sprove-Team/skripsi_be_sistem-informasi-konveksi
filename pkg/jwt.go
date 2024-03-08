package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type JwtC interface {
	CreateToken(refresh bool, claims *Claims, expiresAt time.Time) (string, error)
	ParseToken(refresh bool, token string, claims *Claims) (*jwt.Token, error)
}

type jwtC struct {
	Key        string
	RefreshKey string
}

func NewJwt(key string, refreshKey string) JwtC {
	return &jwtC{key, refreshKey}
}

func (j *jwtC) CreateToken(refresh bool, claims *Claims, expiresAt time.Time) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if refresh {
		return token.SignedString([]byte(j.RefreshKey))
	}
	return token.SignedString([]byte(j.Key))
}

func (j *jwtC) ParseToken(refresh bool, token string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if refresh {
			return []byte(j.RefreshKey), nil
		}
		return []byte(j.Key), nil
	})
}
