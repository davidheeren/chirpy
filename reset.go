package main

import "net/http"

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset"))
	cfg.fileServerHits.Store(0)
}
