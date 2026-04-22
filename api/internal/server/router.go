package server

import (
	"encoding/json"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	registerSwaggerRoutes(mux)

	mux.HandleFunc("POST /api/v1/auth/login", loginHandler)
	mux.HandleFunc("POST /api/v1/articles", createArticleHandler)

	return mux
}

func writeNotImplemented(w http.ResponseWriter, feature, nextStep string) {
	writeJSON(w, http.StatusNotImplemented, notImplementedResponse{
		Message:  feature + " is not implemented yet",
		NextStep: nextStep,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
