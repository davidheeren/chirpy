package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/davidheeren/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	jwtSecret      string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	cfg, err := createConfig()
	if (err != nil) {
		log.Fatal(err)
	}

	server := createServer(cfg)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func createConfig() (*apiConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok {
		return nil, errors.New("DB_URL environment variable must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	return &apiConfig{
		dbQueries: database.New(db),
		jwtSecret: jwtSecret,
	}, nil
}

func createServer(cfg *apiConfig) http.Server {
	serveMux := http.NewServeMux()
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))

	serveMux.HandleFunc("GET /api/healthz", HealthzHandler)
	serveMux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
	serveMux.HandleFunc("GET /api/chirps", cfg.GetChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirpHandler)
	serveMux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	serveMux.HandleFunc("POST /api/login", cfg.LoginUserHandler)
	serveMux.HandleFunc("POST /api/refresh", cfg.RefreshJWTHandler)
	serveMux.HandleFunc("POST /api/revoke", cfg.RevokeRefreshTokenHandler)

	serveMux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	return http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
}
