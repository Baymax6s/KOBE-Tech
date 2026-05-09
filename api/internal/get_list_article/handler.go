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
} // @name server.articleErrorResponse

type TagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleTagJSONResponse

type ArticleListItemJSON struct {
	ID         int64     `json:"id" binding:"required"`
	Title      string    `json:"title" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	UserID     int64     `json:"user_id" binding:"required"`
	Tags       []TagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time `json:"created_at" binding:"required"`
	UpdatedAt  time.Time `json:"updated_at" binding:"required"`
	LikesCount int64     `json:"likes_count" binding:"required"`
} // @name server.articleJSONResponse

type ListArticlesJSONResponse struct {
	Articles []ArticleListItemJSON `json:"articles"`
} // @name server.listArticlesResponse

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
		Articles: make([]ArticleListItemJSON, 0, len(articles)),
	}

	for _, article := range articles {
		response.Articles = append(response.Articles, ArticleListItemJSON{
			ID:         article.ID,
			Title:      article.Title,
			Content:    article.Content,
			UserID:     article.UserID,
			Tags:       newTagJSONs(article.Tags),
			CreatedAt:  article.CreatedAt,
			UpdatedAt:  article.UpdatedAt,
			LikesCount: article.LikesCount,
		})
	}

	return response
}

func newTagJSONs(tags []Tag) []TagJSON {
	response := make([]TagJSON, 0, len(tags))
	for _, tag := range tags {
		response = append(response, TagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}
