package main

import (
	"net/http"
	"sort"

	"github.com/dylst/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")

	var dbChirps []database.Chirp
	var err error

	if authorIdString == "" {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps", err)
			return
		}
	} else {
		authorId, err := uuid.Parse(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse author id", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps", err)
			return
		}
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

	if sortDirection == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
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