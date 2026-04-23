package article

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ArticleJSON struct {
	ID        int64     `json:"id" example:"1"`
	Title     string    `json:"title" example:"First article"`
	Content   string    `json:"content" example:"Article body"`
	UserID    int64     `json:"user_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2026-04-24T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-04-24T00:00:00Z"`
}

type ListArticlesJSONResponse struct {
	Articles []ArticleJSON `json:"articles"`
}

type ErrorJSONResponse struct {
	Message string `json:"message" example:"internal server error"`
}

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListArticles(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.repo == nil {
		log.Printf("list articles: handler is not configured")
		writeJSON(w, http.StatusInternalServerError, ErrorJSONResponse{Message: "internal server error"})
		return
	}

	articles, err := h.repo.List(r.Context())
	if err != nil {
		log.Printf("list articles: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorJSONResponse{Message: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, newListArticlesJSONResponse(articles))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func newListArticlesJSONResponse(articles []Article) ListArticlesJSONResponse {
	response := ListArticlesJSONResponse{
		Articles: make([]ArticleJSON, 0, len(articles)),
	}

	for _, article := range articles {
		response.Articles = append(response.Articles, ArticleJSON{
			ID:        article.ID,
			Title:     article.Title,
			Content:   article.Content,
			UserID:    article.UserID,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		})
	}

	return response
}
