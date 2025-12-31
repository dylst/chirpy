package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dylst/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type response struct{
		User 
		Token string `json:"token"`
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

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second * time.Duration(params.ExpiresInSeconds))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT token", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		},
		Token: token,
	})
}