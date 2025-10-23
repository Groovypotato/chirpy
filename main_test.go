package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/groovypotato/chirpy/internal/auth"
	"github.com/groovypotato/chirpy/internal/database"
	"github.com/joho/godotenv"
)

// ... (your handler definition)

func TestResetHandler(t *testing.T) {
	godotenv.Load()
	var apiCfg apiConfig
	dbURL := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening SQL database: %s", err)
	}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	req, err := http.NewRequest("POST", "/admin/reset", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.resetHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	healthHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	expectedBody := "OK\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}

func TestMetricsHandler(t *testing.T) {
	var apiCfg apiConfig
	req, err := http.NewRequest("GET", "/admin/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.hitsHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	expectedBody := "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited 0 times!</p></body></html>"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}

func TestAppHandler(t *testing.T) {
	var apiCfg apiConfig
	fileHandler := http.FileServer(http.Dir("./"))
	wrapped := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileHandler))

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/app/", nil)

	wrapped.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	expectedBody := "<html>\n  <body>\n    <h1>Welcome to Chirpy</h1>\n  </body>\n</html>"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

}

func TestUsersHandler(t *testing.T) {
	godotenv.Load()
	var apiCfg apiConfig
	ctx := context.Background()
	dbURL := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening SQL database: %s", err)
	}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	payload := map[string]string{"email": "gleasoncr@gmail.com", "password": "password"}
	jsonPayload, _ := json.Marshal(payload)
	reqBody := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest("POST", "/api/users", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.userHandler(rr, req)
	if rr.Code != 201 {
		t.Errorf("Expected status %d, got %d", 201, rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"email":"gleasoncr@gmail.com"`) {
		t.Errorf("email not found in body: %s", body)
	}
	req, err = http.NewRequest("POST", "/admin/reset", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	apiCfg.platform = "dev"
	apiCfg.resetHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reset failed: %d %s", rr.Code, rr.Body.String())
	}

	users, _ := apiCfg.dbQueries.GetUsers(ctx)
	chirps, _ := apiCfg.dbQueries.GetAllChirps(ctx)
	if len(users) != 0 || len(chirps) != 0 {
		t.Fatalf("not cleared: users=%d chirps=%d", len(users), len(chirps))
	}
}

func TestCreateChirpsHandler(t *testing.T) {
	godotenv.Load()
	var apiCfg apiConfig
	ctx := context.Background()
	dbURL := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening SQL database: %s", err)
	}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	user, err := apiCfg.dbQueries.CreateUser(ctx, database.CreateUserParams{
		Email:          "gleasoncr@gmail.com",
		HashedPassword: "x",
	})
	if err != nil {
		t.Error("error creating user")
	}
	payload := map[string]string{"body": "test", "user_id": user.ID.String()}
	jsonPayload, _ := json.Marshal(payload)
	reqBody := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest("POST", "/api/chirps", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.createChirpsHandler(rr, req)
	if rr.Code != 201 {
		t.Errorf("Expected status %d, got %d", 201, rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"body":"test"`) {
		t.Errorf("'test' not found in body: %s", body)
	}
	req, err = http.NewRequest("POST", "/admin/reset", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	apiCfg.platform = "dev"
	apiCfg.resetHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reset failed: %d %s", rr.Code, rr.Body.String())
	}

	users, _ := apiCfg.dbQueries.GetUsers(ctx)
	chirps, _ := apiCfg.dbQueries.GetAllChirps(ctx)
	if len(users) != 0 || len(chirps) != 0 {
		t.Fatalf("not cleared: users=%d chirps=%d", len(users), len(chirps))
	}
}

func TestGetAllChirpsHandler(t *testing.T) {
	godotenv.Load()
	var apiCfg apiConfig
	ctx := context.Background()
	dbURL := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening SQL database: %s", err)
	}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	user, err := apiCfg.dbQueries.CreateUser(ctx, database.CreateUserParams{
		Email:          "gleasoncr@gmail.com",
		HashedPassword: "x",
	})
	if err != nil {
		t.Error("error creating user")
	}
	payload := map[string]string{"body": "test", "user_id": user.ID.String()}
	jsonPayload, _ := json.Marshal(payload)
	reqBody := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest("POST", "/api/chirps", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.createChirpsHandler(rr, req)
	if rr.Code != 201 {
		t.Errorf("Expected status %d, got %d", 201, rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"body":"test"`) {
		t.Errorf("email not found in body: %s", body)
	}
	req, err = http.NewRequest("GET", "/api/chirps", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	apiCfg.getAllChirpsHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	body = rr.Body.String()
	if !strings.Contains(body, `"body":"test"`) {
		t.Errorf("'test' not found in body: %s", body)
	}
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/admin/reset", nil)
	apiCfg.platform = "dev"
	apiCfg.resetHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reset failed: %d %s", rr.Code, rr.Body.String())
	}

	users, _ := apiCfg.dbQueries.GetUsers(ctx)
	chirps, _ := apiCfg.dbQueries.GetAllChirps(ctx)
	if len(users) != 0 || len(chirps) != 0 {
		t.Fatalf("not cleared: users=%d chirps=%d", len(users), len(chirps))
	}
}

func TestGetSingleChirpsHandler(t *testing.T) {
	godotenv.Load()
	var apiCfg apiConfig
	ctx := context.Background()
	dbURL := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Errorf("Error opening database: %s", err.Error())
	}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	user, err := apiCfg.dbQueries.CreateUser(ctx, database.CreateUserParams{
		Email:          "gleasoncr@gmail.com",
		HashedPassword: "x",
	})
	if err != nil {
		t.Errorf("Error creating user: %s", err.Error())
	}
	chirp, err := apiCfg.dbQueries.CreateChirp(ctx, database.CreateChirpParams{
		Body:   "test",
		UserID: user.ID,
	})
	if err != nil {
		t.Errorf("Error creating chirp: %s", err.Error())
	}
	req, err := http.NewRequest("GET", "/api/chirps/"+chirp.ID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("chirpID", chirp.ID.String())
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.getSingleChirpsHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"body":"test"`) {
		t.Errorf(" 'test' not found in body: %s", body)
	}
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/admin/reset", nil)
	apiCfg.platform = "dev"
	apiCfg.resetHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reset failed: %d %s", rr.Code, rr.Body.String())
	}

	users, _ := apiCfg.dbQueries.GetUsers(ctx)
	chirps, _ := apiCfg.dbQueries.GetAllChirps(ctx)
	if len(users) != 0 || len(chirps) != 0 {
		t.Fatalf("not cleared: users=%d chirps=%d", len(users), len(chirps))
	}
}

func TestLoginHandler(t *testing.T) {
	godotenv.Load()
	var apiCfg apiConfig
	ctx := context.Background()
	dbURL := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Errorf("Error opening database: %s", err.Error())
	}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	hpw, _ := auth.HashPassword("x")
	_, err = apiCfg.dbQueries.CreateUser(ctx, database.CreateUserParams{
		Email:          "gleasoncr@gmail.com",
		HashedPassword: hpw,
	})
	if err != nil {
		t.Errorf("Error creating user: %s", err.Error())
	}
	payload := map[string]string{"email": "gleasoncr@gmail.com", "password": "x"}
	jsonPayload, _ := json.Marshal(payload)
	reqBody := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest("POST", "/api/login", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	apiCfg.loginHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"email":"gleasoncr@gmail.com"`) {
		t.Errorf("email not found in body: %s", body)
	}
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/admin/reset", nil)
	apiCfg.platform = "dev"
	apiCfg.resetHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reset failed: %d %s", rr.Code, rr.Body.String())
	}

	users, _ := apiCfg.dbQueries.GetUsers(ctx)
	chirps, _ := apiCfg.dbQueries.GetAllChirps(ctx)
	if len(users) != 0 || len(chirps) != 0 {
		t.Fatalf("not cleared: users=%d chirps=%d", len(users), len(chirps))
	}
}
