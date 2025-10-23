package auth

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeToken(userID uuid.UUID, tokenSecret string) (string, error) {
	currentTime := time.Now().Local().UTC()
	expires := currentTime.Add(1 * time.Hour)
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


func MakeRefreshToken () (string, error) {
	b := make([]byte, 32)
	n, err := rand.Read(b)
	if err != nil {
		return "",err
	}
	binary.BigEndian.PutUint32(b,uint32(n))
	hex := hex.EncodeToString(b)
	return hex, nil
}