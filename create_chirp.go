package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/davidheeren/chirpy/internal/auth"
	"github.com/davidheeren/chirpy/internal/database"
)

func (cfg *apiConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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

	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err = decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid post json", err)
		return
	}

	if len(p.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	// replace profane words
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// ignore new lines I guess
	parts := strings.Split(p.Body, " ")
	for i := 0; i < len(parts); i++ {
		part := strings.ToLower(parts[i])
		if _, ok := profaneWords[part]; ok {
			parts[i] = "****"
		}
	}
	// join parts back (missing previous whitespace other than spaces)
	cleanedBody := strings.Join(parts, " ")

	chirp, err := cfg.dbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: userID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create new chirp in database", err)
		return
	}

	resChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJson(w, http.StatusCreated, resChirp)
}
