package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ValidateChiprHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	p := parameters{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid post json")
		return
	}

	if len(p.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// replace profane words
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	// ignore new lines I guess
	parts := strings.Split(p.Body, " ")
	for i := 0; i < len(parts); i++ {
		part := strings.ToLower(parts[i])
		if _, ok := profaneWords[part]; ok {
			parts[i] = "****"
		}
	}
	res := returnVals{
		// join parts back (missing previous whitespace other than spaces)
		CleanedBody: strings.Join(parts, " "),
	}

	respondWithJson(w, http.StatusOK, res)
}
