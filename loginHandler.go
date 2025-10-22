package main

import (
	"encoding/json"
	"net/http"

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
	} else {
		respondWithJSON(w, 200, userResp{ID: vuser.ID, CreatedAt: vuser.CreatedAt, UpdatedAt: vuser.UpdatedAt, Email: vuser.Email})
	}
}
