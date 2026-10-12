package main

import (
	"net/http"
	"sort"
	"strings"

	"github.com/davidheeren/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortMethod := strings.ToLower(r.URL.Query().Get("sort"))
	if sortMethod != "asc" && sortMethod != "desc" {
		sortMethod = "asc"
	}

	var chirps []database.Chirp
	var err error
	// get chirps only by user id if author_id is not empty
	// otherwise get all chirps
	if authorID != "" {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps in database", err)
			return
		}
		chirps, err = cfg.dbQueries.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps in database", err)
			return
		}
	} else {
		chirps, err = cfg.dbQueries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps in database", err)
			return
		}
	}

	if sortMethod == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}

	resChirps := make([]Chirp, len(chirps))

	for i, c := range chirps {
		resChirps[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	respondWithJson(w, http.StatusOK, resChirps)
}
