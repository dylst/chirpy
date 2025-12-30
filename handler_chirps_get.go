package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps", err)
	}

	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	respondWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid Chirp ID", err)
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to fetch chirp", err)
		return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	})
}