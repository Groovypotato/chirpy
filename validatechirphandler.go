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
	decoder := json.NewDecoder(r.Body)
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
		badWordFilterHandler(w,200,params.Body)
		return
	}
}