package article

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.articleErrorResponse

type AuthorJSON struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} // @name server.articleAuthorJSONResponse

type ArticleJSON struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Author    AuthorJSON `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
} // @name server.getArticleJSONResponse

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/articles/:article_id", h.getArticleHandler)
}

// getArticleHandler godoc
//
//	@Summary		Get article
//	@Description	Get article detail API.
//	@Tags			articles
//	@Produce		json
//	@Param			article_id	path		int	true	"Article ID"
//	@Success		200			{object}	ArticleJSON
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/articles/{article_id} [get]
func (h *Handler) getArticleHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid article_id"})
		return
	}

	response, err := h.GetArticle(c.Request.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, errArticleNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get article"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

var errArticleNotFound = errors.New("article not found")

func (h *Handler) GetArticle(ctx context.Context, articleID int64) (ArticleJSON, error) {
	if h == nil || h.repo == nil {
		return ArticleJSON{}, errors.New("get article handler is not configured")
	}

	article, err := h.repo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ArticleJSON{}, errArticleNotFound
		}
		return ArticleJSON{}, err
	}

	return ArticleJSON{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		Author: AuthorJSON{
			ID:   article.Author.ID,
			Name: article.Author.Name,
		},
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}
