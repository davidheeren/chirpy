package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) PolkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err := decoder.Decode(&p)
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
