package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid post json", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), p.Email)
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
