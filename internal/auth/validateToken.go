package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}
	u, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return u, nil
}
