package main

import (
	"encoding/json"
	"net/http"

	"github.com/davidheeren/chirpy/internal/auth"
)

func (cfg *apiConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
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
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	matchPassword, err := auth.CheckPasswordHash(p.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if !matchPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	rv := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJson(w, http.StatusOK, rv)
}
