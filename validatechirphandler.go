package main

import (
	"encoding/json"
	"net/http"
)

func validatechirpHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST"{
		respondWithError(w,405,"method not allowed")
	}
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValid struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w,400,"Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w,400,"Chirp is too long")
		return
	}else {
		respondWithJSON(w,200, returnValid{Valid: true})
		return
	}
}