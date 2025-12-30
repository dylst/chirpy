package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dylst/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body 	  string 	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	/* 
	define the params 
	init the struct 
	define the decoder 
	decode the body into params

	validate the body
	*/
	type parameters struct{
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode params", err)
	}

	const maxChirpLength = 140 
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {}, 
		"sharbert": {}, 
		"fornax": {},
	}
	cleanedBodyResponse := getCleanedBody(params.Body, profaneWords)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBodyResponse,
		UserID: params.UserID,
	})

	respondWithJson(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func getCleanedBody(body string, profaneWords map[string]struct{}) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}