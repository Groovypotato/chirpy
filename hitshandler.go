package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) hitsHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d\n",hits)
}

func (cfg *apiConfig) resethitsHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}