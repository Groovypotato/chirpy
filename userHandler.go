package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	type userEmail struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	uemail := userEmail{}
	err := decoder.Decode(&uemail)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	user, err := cfg.dbQueries.GetUser(r.Context(), uemail.Email)
	if err != nil {
		user, err = cfg.dbQueries.CreateUser(r.Context(), uemail.Email)
		if err != nil {
			respondWithError(w, 400, "something went wrong")
			return
		}
		respondWithJSON(w, 201, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
	} else {
		respondWithError(w, 409, "user already exists")
		return
	}

}
