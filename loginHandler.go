package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/groovypotato/chirpy/internal/auth"
	"github.com/groovypotato/chirpy/internal/database"
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
	rexpire := time.Now().AddDate(0,0,60)
	token, err := auth.MakeToken(vuser.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	rtoken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "unable to make refresh token")
	}
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(),database.CreateRefreshTokenParams{
		Token: rtoken,
		UserID: vuser.ID,
		ExpiresAt: rexpire,
	})
	if err != nil {
		respondWithError(w, 500, "unable to insert refresh token in db")
		return
	}
	respondWithJSON(w, 200, userResp{ID: vuser.ID, CreatedAt: vuser.CreatedAt, UpdatedAt: vuser.UpdatedAt, Email: vuser.Email, Token: token, RefreshToken: rtoken})
}
