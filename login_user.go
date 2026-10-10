package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/davidheeren/chirpy/internal/auth"
	"github.com/davidheeren/chirpy/internal/database"
)

type loginUser struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid post json", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), p.Email)
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

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create jwt", err)
		return
	}

	tokenExpiresAt := time.Now().UTC().UTC().Add(time.Hour * 24 * 60) // 60 days
	rft, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: tokenExpiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create refresh token", err)
		return
	}

	rv := loginUser{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        jwt,
		RefreshToken: rft.Token,
	}

	respondWithJson(w, http.StatusOK, rv)
}
