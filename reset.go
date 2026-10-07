package main

import "net/http"

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	err := cfg.dbQueries.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not reset the database")
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset"))
	cfg.fileServerHits.Store(0)
}
