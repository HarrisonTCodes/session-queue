package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSecret = []byte("secret")

func CreateToken(position int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": position,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecret)
}
