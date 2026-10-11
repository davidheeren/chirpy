package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/davidheeren/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) PolkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate api key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "could not validate api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err = decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid post json", err)
		return
	}

	if p.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "not a valid polka webook", errors.New("not a valid polka webook"))
		return
	}

	_, err = cfg.dbQueries.UpdateUserChirpyRed(r.Context(), p.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not update user's chirpy red", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
