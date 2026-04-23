package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
)

type apiServer struct {
	articleHandler *article.Handler
}

func NewHandler(db *sql.DB) http.Handler {
	server := &apiServer{
		articleHandler: article.NewHandler(article.NewRepository(db)),
	}

	mux := http.NewServeMux()

	registerSwaggerRoutes(mux)

	mux.HandleFunc("POST /api/auth/login", server.loginHandler)
	mux.HandleFunc("GET /api/articles", server.listArticlesHandler)
	mux.HandleFunc("POST /api/articles", server.createArticleHandler)

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
