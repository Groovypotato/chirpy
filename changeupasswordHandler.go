package main

import (
	"encoding/json"
	"net/http"

	"github.com/groovypotato/chirpy/internal/auth"
	"github.com/groovypotato/chirpy/internal/database"
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
	uid, err := auth.ValidateJWT(jwt, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: bad jwt")
	}
	decoder := json.NewDecoder(r.Body)
	uinput := userInput{}
	err = decoder.Decode(&uinput)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: issue decoding")
		return
	}
	hpass, err := auth.HashPassword(uinput.Password)
	if err != nil {
		respondWithError(w, 400, "something went wrong: unable to hash password")
		return
	}
	err = cfg.dbQueries.ChangePassword(r.Context(), database.ChangePasswordParams{
		HashedPassword: hpass,
		ID:             uid,
		Email:          uinput.Email,
	})

	if err != nil {
		respondWithError(w, 401, "error chaning password in db")
	}
	newJWT, err := auth.MakeToken(uid, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	vuser, err := cfg.dbQueries.GetUser(r.Context(), uinput.Email)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	rtoken, err := cfg.dbQueries.GetUserRefreshToken(r.Context(), uid)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 200, userResp{
		ID:           uid,
		CreatedAt:    vuser.CreatedAt,
		UpdatedAt:    vuser.UpdatedAt,
		Email:        vuser.Email,
		Token:        newJWT,
		RefreshToken: rtoken.Token,
	})

}
