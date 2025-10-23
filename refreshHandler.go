package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/groovypotato/chirpy/internal/auth"
	"github.com/groovypotato/chirpy/internal/database"
)


func (cfg *apiConfig) validateRefreshToken(r *http.Request) (database.RefreshToken, error) {
	currentTime := time.Now()
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return database.RefreshToken{}, err
	}
	vtoken, err := cfg.dbQueries.GetRefreshToken(r.Context(),token)
	if err != nil {
		return database.RefreshToken{}, err
	}
	if currentTime.After(vtoken.ExpiresAt) {
		return database.RefreshToken{}, errors.New("unauthorized: token has expired")
	}
	return vtoken, nil
}


func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	vtoken, err := cfg.validateRefreshToken(r)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	naToken, err := auth.MakeToken(vtoken.UserID,cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 400, "someting went wrong making new access token")
		return
	}
	respondWithJSON(w, 200,Token{TOKEN: naToken})
}


 func (cfg *apiConfig) refreshRevokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	vtoken, err := cfg.validateRefreshToken(r)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(),vtoken.Token)
	if err != nil {
		respondWithError(w,401,"something went wrong revoking token")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(204)
 }