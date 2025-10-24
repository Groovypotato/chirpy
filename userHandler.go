package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/groovypotato/chirpy/internal/auth"
	"github.com/groovypotato/chirpy/internal/database"
)

func (cfg *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	var user database.User
	decoder := json.NewDecoder(r.Body)
	uinput := userInput{}
	err := decoder.Decode(&uinput)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: issue decoding")
		return
	}
	_, err = cfg.dbQueries.GetUser(r.Context(), uinput.Email)
	if err != nil {
		hashedPassword, err := auth.HashPassword(uinput.Password)
		if err != nil {
			respondWithError(w, 400, "something went wrong: issue hashing password")
			return
		}
		params := database.CreateUserParams{
			Email:          uinput.Email,
			HashedPassword: hashedPassword,
		}
		user, err = cfg.dbQueries.CreateUser(r.Context(), params)
		if err != nil {
			respondWithError(w, 400, "something went wrong:issue creating user")
			return
		}
		respondWithJSON(w, 201, userResp{ID: user.ID, 
			CreatedAt: user.CreatedAt, 
			UpdatedAt: user.UpdatedAt, 
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		})
	} else {
		respondWithError(w, 409, "user already exists")
		return
	}
}

func (cfg *apiConfig) upgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	decoder := json.NewDecoder(r.Body)
	whook := WebhookEvent{}
	err := decoder.Decode(&whook)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: issue decoding")
		return
	}
	if whook.Event != "user.upgraded" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.WriteHeader(204)
	}
	uid, err := uuid.Parse(whook.Data.UserID)
	if err != nil {
		respondWithError(w,400,err.Error())
		return
	}
	err = cfg.dbQueries.UpgradeUserChirpyRed(r.Context(), uid)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			respondWithError(w,404,"user not found")
			return
		}
		respondWithError(w,400,"error upgrading user")
	}
	w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.WriteHeader(204)
}
