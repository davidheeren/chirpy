package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/davidheeren/chirpy/internal/auth"
	"github.com/davidheeren/chirpy/internal/database"
)

func (cfg *apiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	if len(p.Password) < 4 {
		respondWithError(w, http.StatusInternalServerError, "password must be at leas 4 characters", errors.New("password cannot be less than 4 characters"))
		return
	}

	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email:          p.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
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
