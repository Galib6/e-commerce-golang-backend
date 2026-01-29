package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
)

type JWTClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"` // ["user", "admin"]
	jwt.RegisteredClaims
}

func CreateToken(claims JWTClaims) (string, error) {
	env := config.LoadEnv()
	secret := []byte(env.JWTSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*JWTClaims, error) {
	env := config.LoadEnv()
	secret := []byte(env.JWTSecret)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
