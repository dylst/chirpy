package main

import (
	"net/http"

	"github.com/dylst/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
/*
Add a new DELETE /api/chirps/{chirpID} route to your server that deletes a chirp from the database by its id.
This is an authenticated endpoint, so be sure to check the token in the header. Only allow the deletion of a chirp if the user is the author of the chirp.
If they are not, return a 403 status code.
If the chirp is deleted successfully, return a 204 status code.
If the chirp is not found, return a 404 status code.
*/
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing bearer token", err)
		return
	}

	chirpId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse chirp id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
		return
	}

	if userId != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "User if not the author of the chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}