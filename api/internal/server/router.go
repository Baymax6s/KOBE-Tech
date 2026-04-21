package server

import (
	"encoding/json"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/login", notImplemented("login", "internal/auth/{handler,service,repository}.go"))
	mux.HandleFunc("POST /api/v1/articles", notImplemented("create article", "internal/article/{handler,service,repository}.go"))

	return mux
}

func notImplemented(feature, nextStep string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotImplemented, map[string]string{
			"message":   feature + " is not implemented yet",
			"next_step": nextStep,
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
