package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/groovypotato/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening SQL database: %s", err)
	}
	var apiCfg apiConfig
	apiCfg.dbQueries = database.New(db)
	mux := http.NewServeMux()
	fileHandler := http.FileServer(http.Dir("./"))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileHandler)))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validatechirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.userHandler)
	srv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("error when starting server:%v ", err)
	}
}
