package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthorJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleAuthorJSONResponse

type GetArticleJSONResponse struct {
	ID         int64            `json:"id" binding:"required"`
	Title      string           `json:"title" binding:"required"`
	Content    string           `json:"content" binding:"required"`
	Author     AuthorJSON       `json:"author" binding:"required"`
	Tags       []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time        `json:"created_at" binding:"required"`
	UpdatedAt  time.Time        `json:"updated_at" binding:"required"`
	LikesCount int64            `json:"likes_count" binding:"required"`
} // @name server.getArticleJSONResponse

// getArticleHandler godoc
//
//	@Summary		Get article
//	@Description	Get article detail API.
//	@Tags			articles
//	@Produce		json
//	@Param			article_id	path		int	true	"Article ID"
//	@Success		200			{object}	GetArticleJSONResponse
//	@Failure		400			{object}	ArticleErrorResponse
//	@Failure		404			{object}	ArticleErrorResponse
//	@Failure		500			{object}	ArticleErrorResponse
//	@Router			/api/articles/{article_id} [get]
func (h *Handler) getArticleHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ArticleErrorResponse{Message: "invalid article_id"})
		return
	}

	response, err := h.GetArticle(c.Request.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, errArticleNotFound):
			c.JSON(http.StatusNotFound, ArticleErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ArticleErrorResponse{Message: "failed to get article"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

var errArticleNotFound = errors.New("article not found")

func (h *Handler) GetArticle(ctx context.Context, articleID int64) (GetArticleJSONResponse, error) {
	if h == nil || h.repo == nil {
		return GetArticleJSONResponse{}, errors.New("get article handler is not configured")
	}

	item, err := h.repo.FindArticleByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetArticleJSONResponse{}, errArticleNotFound
		}
		return GetArticleJSONResponse{}, err
	}

	return GetArticleJSONResponse{
		ID:      item.ID,
		Title:   item.Title,
		Content: item.Content,
		Author: AuthorJSON{
			ID:   item.Author.ID,
			Name: item.Author.Name,
		},
		Tags:       newArticleTagJSONs(item.Tags),
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
		LikesCount: item.LikesCount,
	}, nil
}
