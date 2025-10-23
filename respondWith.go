package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	errorMsg := returnError {
			Error: msg,
		}
		dat, err := json.Marshal(errorMsg)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(code)
			fmt.Fprintf(w, "error: %s\n", msg)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.WriteHeader(code)
		w.Write(dat)
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	validMsg := payload
		dat, err := json.Marshal(validMsg)
		if err != nil {
			respondWithError(w,400,"Something went wrong")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.WriteHeader(code)
		w.Write(dat)
}