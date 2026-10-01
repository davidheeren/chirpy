package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	serveMux := http.NewServeMux()
	cfg := &apiConfig{}

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))
	serveMux.HandleFunc("GET /api/healthz", HealthzHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
