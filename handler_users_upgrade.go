package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Event string `json:"event"`
		Data  struct{
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
	}

	/*
	If the event is anything other than user.upgraded, the endpoint should immediately respond with a 204 status code - we don't care about any other events. (done)
If the event is user.upgraded, then it should update the user in the database, and mark that they are a Chirpy Red member.
If the user is upgraded successfully, the endpoint should respond with a 204 status code and an empty response body. If the user can't be found, the endpoint should respond with a 404 status code.
	*/
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Update user 
	_, err = cfg.db.UpgradeUser(r.Context(), params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, struct{}{})
}