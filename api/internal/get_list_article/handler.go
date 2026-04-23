package article

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

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

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/articles", h.listArticlesHandler)
}

// listArticlesHandler godoc
//
//	@Summary		List articles
//	@Description	Get article list API.
//	@Tags			articles
//	@Produce		json
//	@Success		200	{object}	ListArticlesJSONResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/articles [get]
func (h *Handler) listArticlesHandler(c *gin.Context) {
	response, err := h.ListArticles(c.Request.Context())
	if err != nil {
		log.Printf("list articles: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
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
