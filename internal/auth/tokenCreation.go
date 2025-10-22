package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func makeToken(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currentTime := time.Now().Local().UTC()
	expires := currentTime.Add(expiresIn)
	tokenSecretByte := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{Issuer: "chirpy",
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expires),
			Subject:   userID.String()})
	sstring, err := token.SignedString(tokenSecretByte)
	if err != nil {
		return "", err
	}
	return sstring, nil
}
