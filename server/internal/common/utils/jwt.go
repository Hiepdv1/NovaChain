package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

func SignJWT[T any](secret []byte, data T, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := JWTPayload[T]{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(secret)
}

func VerifyJWT[T any](secret []byte, tokenStr string) (*JWTPayload[T], error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTPayload[T]{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token")
			}
			return secret, nil
		},
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token is expired")
		}
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTPayload[T])
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
