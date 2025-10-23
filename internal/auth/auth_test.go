package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenCreation(t *testing.T) {
	uid := uuid.New()
	d := 30 * time.Hour
	secret := "4$KGXaO8+LNuOIKjqQoi/PIg$GOk8JWFN4zQXqs+bdrixjOgF7lt8dsjwYI0U/Y5E4+w"
	result, err := MakeToken(uid, secret, d)
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

func TestGetbearer(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Custom-Header", "my-value")
	req.Header.Set("Authorization", "Bearer my-token")
	token, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("there was an issue getting the token: %s", err.Error())
	}
	if token != "my-token" {
		t.Errorf("token expected: 'my-token' token receieved:'%s' ", token)
	}

}
