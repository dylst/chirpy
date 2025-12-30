package main

import (
	"encoding/json"
	"net/http"

	"github.com/dylst/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode params", err)
		return
	}

	user, err := cfg.db.GetUserForEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch user by email", err)
		return
	}
	hashedPassword := user.HashedPassword

	match, err := auth.CheckPasswordHash(params.Password, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to validate password", err)
		return
	}

	if !match{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}