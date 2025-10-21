package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respondWithError(w, 405, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	httpformattedstring := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	byteSlice := []byte(httpformattedstring)
	w.Write(byteSlice)
}
