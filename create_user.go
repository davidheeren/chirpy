package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/davidheeren/chirpy/internal/auth"
	"github.com/davidheeren/chirpy/internal/database"
)

func (cfg *apiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if len(p.Password) < 4 {
		respondWithError(w, http.StatusInternalServerError, "password must be at leas 4 characters", errors.New("password cannot be less than 4 characters"))
		return
	}

	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: p.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create new user", err)
		return
	}

	rv := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJson(w, http.StatusCreated, rv)
}
