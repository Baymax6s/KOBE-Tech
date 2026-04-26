package article

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type createArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
} // @name server.createArticleRequest

type ArticleJSON struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name server.articleJSONResponse

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.articleErrorResponse

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/articles", h.createArticleHandler)
}

// createArticleHandler godoc
//
//	@Summary		Create article
//	@Description	Create article API.
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createArticleRequest	true	"Create article request"
//	@Success		201		{object}	ArticleJSON
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/articles [post]
func (h *Handler) createArticleHandler(c *gin.Context) {
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	userID := auth.MustUserID(c)

	response, err := h.CreateArticle(c.Request.Context(), userID, req.Title, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidRequest):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to create article"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

var errInvalidRequest = errors.New("title and content are required")

func (h *Handler) CreateArticle(ctx context.Context, userID int64, title, content string) (ArticleJSON, error) {
	if h == nil || h.repo == nil {
		return ArticleJSON{}, errors.New("post article handler is not configured")
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" || content == "" {
		return ArticleJSON{}, errInvalidRequest
	}

	article, err := h.repo.Create(ctx, title, content, userID)
	if err != nil {
		return ArticleJSON{}, err
	}

	return ArticleJSON{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		UserID:    article.UserID,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}
