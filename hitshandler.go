package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) hitsHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	httpformattedstring := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	byteSlice := []byte(httpformattedstring)
	w.Write(byteSlice)
}

func (cfg *apiConfig) resethitsHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}