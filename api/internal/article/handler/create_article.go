package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type createArticleRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	TagNames []string `json:"tags" minLength:"1" maxLength:"10"`
} // @name server.createArticleRequest

type CreateArticleJSONResponse struct {
	ID        int64            `json:"id" binding:"required"`
	Title     string           `json:"title" binding:"required"`
	Content   string           `json:"content" binding:"required"`
	UserID    int64            `json:"user_id" binding:"required"`
	Tags      []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt time.Time        `json:"created_at" binding:"required"`
	UpdatedAt time.Time        `json:"updated_at" binding:"required"`
} // @name server.createArticleResponse

// createArticleHandler godoc
//
//	@Summary		Create article
//	@Description	Create article API.
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createArticleRequest	true	"Create article request"
//	@Success		201		{object}	CreateArticleJSONResponse
//	@Failure		400		{object}	ArticleErrorResponse
//	@Failure		401		{object}	ArticleErrorResponse
//	@Failure		500		{object}	ArticleErrorResponse
//	@Security		BearerAuth
//	@Router			/api/articles [post]
func (h *Handler) createArticleHandler(c *gin.Context) {
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ArticleErrorResponse{Message: "invalid request body"})
		return
	}

	userID := auth.MustUserID(c)

	response, err := h.CreateArticle(c.Request.Context(), userID, req.Title, req.Content, req.TagNames)
	if err != nil {
		switch {
		case errors.Is(err, article.ErrInvalidArticle), errors.Is(err, article.ErrInvalidTagName):
			c.JSON(http.StatusBadRequest, ArticleErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ArticleErrorResponse{Message: "failed to create article"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) CreateArticle(ctx context.Context, userID int64, title, content string, tagNames []string) (CreateArticleJSONResponse, error) {
	if h == nil || h.repo == nil {
		return CreateArticleJSONResponse{}, errors.New("post article handler is not configured")
	}

	title, content, normalizedTagNames, err := article.NormalizeCreateInput(title, content, tagNames)
	if err != nil {
		return CreateArticleJSONResponse{}, err
	}

	item, err := h.repo.CreateArticle(ctx, title, content, userID, normalizedTagNames)
	if err != nil {
		return CreateArticleJSONResponse{}, err
	}

	return CreateArticleJSONResponse{
		ID:        item.ID,
		Title:     item.Title,
		Content:   item.Content,
		UserID:    item.UserID,
		Tags:      newArticleTagJSONs(item.Tags),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
