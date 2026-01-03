package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dylst/chirpy/internal/auth"
	"github.com/dylst/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct{
		User 
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second * 3600)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken() 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	err = cfg.createRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
		},
		Token: token,
		RefreshToken: refreshToken,
	})
}