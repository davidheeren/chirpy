package main

import (
	"net/http"

	"github.com/davidheeren/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("id")
	chirpID, err := uuid.Parse(pathID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse chrip id from path", err)
		return
	}

	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate user. please login", err)
		return
	}

	userID, err := auth.ValidateJWT(jwt, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate user. please login", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not find chirp in database", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "could not validate user. please login", err)
		return

	}
	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "could not delete chirp in database", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
