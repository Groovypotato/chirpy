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
	newType := token.Claims.(*jwt.RegisteredClaims)
	newU, err := uuid.Parse(newType.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return newU, nil
}
