package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("id")
	chirpID, err := uuid.Parse(pathID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse chrip id from path", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not get chirp in database", err)
		return
	}

	resChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJson(w, http.StatusOK, resChirp)
}
