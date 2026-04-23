package article

import (
	"context"
	"errors"
	"time"
)

type ArticleJSON struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListArticlesJSONResponse struct {
	Articles []ArticleJSON `json:"articles"`
}

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListArticles(ctx context.Context) (ListArticlesJSONResponse, error) {
	if h == nil || h.repo == nil {
		return ListArticlesJSONResponse{}, errors.New("article handler is not configured")
	}

	articles, err := h.repo.List(ctx)
	if err != nil {
		return ListArticlesJSONResponse{}, err
	}

	return newListArticlesJSONResponse(articles), nil
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
