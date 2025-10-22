package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenCreation(t *testing.T) {
	uid := uuid.New()
	d := 30 * time.Hour
	secret := "4$KGXaO8+LNuOIKjqQoi/PIg$GOk8JWFN4zQXqs+bdrixjOgF7lt8dsjwYI0U/Y5E4+w"
	result, err := makeToken(uid, secret, d)
	if err != nil {
		t.Errorf("error making token: %s", err.Error())
	}
	if len(result) == 0 {
		t.Error("token length 0")
	}
	_, err = ValidateJWT(result, secret)
	if err != nil {
		t.Errorf("error validating token: %s", err.Error())
	}

}
