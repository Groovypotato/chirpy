package main

import (
	"encoding/json"
	"net/http"

	"github.com/groovypotato/chirpy/internal/auth"
)

func (cfg *apiConfig) changeUPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: no bearer token")
		return
	}
	uid, err :=auth.ValidateJWT(jwt,cfg.jwtSecret)
	if err != nil {
		respondWithError(w,400,"Unauthorized: bad jwt")
	}
	decoder := json.NewDecoder(r.Body)
	uinput := userInput{}
	err = decoder.Decode(&uinput)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: issue decoding")
		return
	}
	
}