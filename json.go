package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type returnVals struct {
		Error string `json:"error"`
	}
	if err != nil {
		log.Println(err)
	}
	res := returnVals{Error: msg}
	respondWithJson(w, code, res)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error marshalling JSON: %s", err)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
