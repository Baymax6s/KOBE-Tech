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
	repo      *Repository
	validator *auth.Validator
}

func NewHandler(repo *Repository, validator *auth.Validator) *Handler {
	return &Handler{repo: repo, validator: validator}
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
//	@Router			/api/articles [post]
func (h *Handler) createArticleHandler(c *gin.Context) {
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	response, err := h.CreateArticle(c.Request.Context(), c.GetHeader("Authorization"), req.Title, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidAuthorization), errors.Is(err, errInvalidToken):
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		case errors.Is(err, errInvalidRequest):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to create article"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

var (
	errInvalidAuthorization = errors.New("invalid authorization header")
	errInvalidToken         = errors.New("invalid token")
	errInvalidRequest       = errors.New("title and content are required")
)

func (h *Handler) CreateArticle(ctx context.Context, authorization, title, content string) (ArticleJSON, error) {
	if h == nil || h.repo == nil || h.validator == nil {
		return ArticleJSON{}, errors.New("post article handler is not configured")
	}

	parts := strings.Fields(strings.TrimSpace(authorization))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ArticleJSON{}, errInvalidAuthorization
	}

	userID, err := h.validator.ValidateToken(parts[1])
	if err != nil {
		return ArticleJSON{}, errInvalidToken
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
