package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/davidheeren/chirpy/internal/auth"
)

func (cfg *apiConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid post json", err)
		return
	}

	user, err := cfg.dbQueries.GeUserByEmail(r.Context(), p.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	matchPassword, err := auth.CheckPasswordHash(p.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if !matchPassword {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	expiresIn := time.Hour
	if p.ExpiresInSeconds != 0 {
		expiresIn = time.Second * time.Duration(p.ExpiresInSeconds)
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create jwt", err)
		return
	}

	rv := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwt,
	}

	respondWithJson(w, http.StatusOK, rv)
}
