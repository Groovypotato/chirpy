package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/groovypotato/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
}

type userInput struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type userResp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.jwtSecret = os.Getenv("JWT_SECRET")
	mux := http.NewServeMux()
	fileHandler := http.FileServer(http.Dir("./"))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileHandler)))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.userHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpsHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getSingleChirpsHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	srv := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("error when starting server:%v ", err)
	}
}
