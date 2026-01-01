package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dylst/chirpy/internal/auth"
	"github.com/dylst/chirpy/internal/database"
)

func (cfg *apiConfig) createRefreshToken(ctx context.Context, params database.CreateRefreshTokenParams) error {
	_, err := cfg.db.CreateRefreshToken(ctx, params)
	return err
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get bearer token", err)
		return
	}

	token, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get user from token", err)
		return
	}

	if time.Now().After(token.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Expired token", err)
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token revoked", nil)
    	return
	}

	access_token, err := auth.MakeJWT(token.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT token", err)
		return
	}

	type response struct{
		Token string `json:"token"`
	}

	respondWithJson(w, http.StatusOK, response{
		Token: access_token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get bearer token", err)
		return
	}

	err = cfg.db.UpdateRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}