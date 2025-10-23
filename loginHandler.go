package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/groovypotato/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	decoder := json.NewDecoder(r.Body)
	uinput := userInput{}
	err := decoder.Decode(&uinput)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	vuser, err := cfg.dbQueries.GetUser(r.Context(), uinput.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	ok, err := auth.CheckPasswordHash(uinput.Password, vuser.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if uinput.ExpiresInSeconds <= 0 {
		uinput.ExpiresInSeconds = 3600
	} else if uinput.ExpiresInSeconds > 3600 {
		uinput.ExpiresInSeconds = 3600
	}
	dur := time.Duration(uinput.ExpiresInSeconds) * time.Second
	token, err := auth.MakeToken(vuser.ID, cfg.jwtSecret, dur)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	respondWithJSON(w, 200, userResp{ID: vuser.ID, CreatedAt: vuser.CreatedAt, UpdatedAt: vuser.UpdatedAt, Email: vuser.Email, Token: token})
}
