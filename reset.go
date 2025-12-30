package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	cfg.fileServerHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete users", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
