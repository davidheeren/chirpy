package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/davidheeren/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serveMux := http.NewServeMux()
	cfg := &apiConfig{
		dbQueries: database.New(db),
	}

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))

	serveMux.HandleFunc("GET /api/healthz", HealthzHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", ValidateChiprHandler)

	serveMux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
