package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error::JSON marshal::%v::%v\n", err, payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code >= 500 {
		log.Printf("respondWithError::info::code%v::message%v\n", code, message)
	}

	type errResponse struct {
		Error string `json:"error"` // Maps to "error" in JSON instead of "Error"
	}

	respondWithJSON(w, code, errResponse{
		Error: message,
	})
}
