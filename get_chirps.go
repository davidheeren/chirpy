package main

import (
	"net/http"
)

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not chirps in database", err)
		return
	}

	resChirps := make([]Chirp, len(chirps))

	for i, c := range chirps {
		resChirps[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	respondWithJson(w, http.StatusOK, resChirps)
}
