package main

import (
	"encoding/json"
	"net/http"

	"github.com/dylst/chirpy/internal/auth"
	"github.com/dylst/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing bearer token", err)
		return
	}

	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}
	respondWithJson(w, http.StatusOK, User{
		ID: userId, 
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}