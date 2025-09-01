package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(position int64, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": position,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 3).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString(secret)
}

func ValidateToken(tknString string, secret []byte) (int64, error) {
	tkn, err := jwt.Parse(tknString, func(tkn *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return 0, errors.New(err.Error())
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("claims are not MapClaims")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("sub claim is not a number")
	}

	return int64(sub), nil
}
