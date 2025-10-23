package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/groovypotato/chirpy/internal/auth"
	"github.com/groovypotato/chirpy/internal/database"
)

type Chirp struct {
	Body string `json:"body"`
}

type VChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	uid, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 400, "Something went wrong with the decoder")
		return
	}
	valid, code, post := validateChirpHandler(chirp.Body)
	if !valid {
		respondWithError(w, code, post)
		return
	}

	parms := database.CreateChirpParams{
		Body:   post,
		UserID: uid,
	}
	vchirp, err := cfg.dbQueries.CreateChirp(r.Context(), parms)
	if err != nil {
		errorString := err.Error()
		respondWithError(w, 400, errorString)
		return
	}
	respondWithJSON(w, 201, VChirp{ID: vchirp.ID, CreatedAt: vchirp.CreatedAt, UpdatedAt: vchirp.UpdatedAt, Body: vchirp.Body, UserId: vchirp.UserID})
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	allChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "something went wrong")
	}
	resp := make([]VChirp, 0, len(allChirps))
	for _, c := range allChirps {
		resp = append(resp, VChirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		})
	}

	// prefer your helper for consistency
	respondWithJSON(w, 200, resp)

}

func (cfg *apiConfig) getSingleChirpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	id, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, 400, "something went wrong")
		return
	}
	vchirp, err := cfg.dbQueries.GetOneChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}

	respondWithJSON(w, 200, VChirp{ID: vchirp.ID, CreatedAt: vchirp.CreatedAt, UpdatedAt: vchirp.UpdatedAt, Body: vchirp.Body, UserId: vchirp.UserID})
}
