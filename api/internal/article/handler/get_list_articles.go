package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/gin-gonic/gin"
)

type ArticleErrorResponse struct {
	Message string `json:"message"`
} // @name server.articleErrorResponse

type ArticleTagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleTagJSONResponse

type ArticleListItemJSON struct {
	ID         int64            `json:"id" binding:"required"`
	Title      string           `json:"title" binding:"required"`
	Content    string           `json:"content" binding:"required"`
	UserID     int64            `json:"user_id" binding:"required"`
	Tags       []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time        `json:"created_at" binding:"required"`
	UpdatedAt  time.Time        `json:"updated_at" binding:"required"`
	LikesCount int64            `json:"likes_count" binding:"required"`
} // @name server.articleJSONResponse

type ListArticlesJSONResponse struct {
	Articles []ArticleListItemJSON `json:"articles"`
} // @name server.listArticlesResponse

// listArticlesHandler godoc
//
//	@Summary		List articles
//	@Description	Get article list API.
//	@Tags			articles
//	@Produce		json
//	@Success		200	{object}	ListArticlesJSONResponse
//	@Failure		500	{object}	ArticleErrorResponse
//	@Router			/api/articles [get]
func (h *Handler) listArticlesHandler(c *gin.Context) {
	response, err := h.ListArticles(c.Request.Context())
	if err != nil {
		log.Printf("list articles: %v", err)
		c.JSON(http.StatusInternalServerError, ArticleErrorResponse{
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

	articles, err := h.repo.ListArticles(ctx)
	if err != nil {
		return ListArticlesJSONResponse{}, err
	}

	return newListArticlesJSONResponse(articles), nil
}

func newListArticlesJSONResponse(articles []article.Article) ListArticlesJSONResponse {
	response := ListArticlesJSONResponse{
		Articles: make([]ArticleListItemJSON, 0, len(articles)),
	}

	for _, item := range articles {
		response.Articles = append(response.Articles, ArticleListItemJSON{
			ID:         item.ID,
			Title:      item.Title,
			Content:    item.Content,
			UserID:     item.UserID,
			Tags:       newArticleTagJSONs(item.Tags),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			LikesCount: item.LikesCount,
		})
	}

	return response
}

func newArticleTagJSONs(tags []article.Tag) []ArticleTagJSON {
	response := make([]ArticleTagJSON, 0, len(tags))
	for _, tag := range tags {
		response = append(response, ArticleTagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}
