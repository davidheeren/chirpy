package main

import (
	"net/http"
	"time"

	"github.com/davidheeren/chirpy/internal/auth"
)

type refreshedToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) RefreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate user. please login", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token not valid", err)
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create jwt", err)
		return
	}

	rv := refreshedToken{
		Token: jwt,
	}

	respondWithJson(w, http.StatusOK, rv)
}
