package like

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.likeErrorResponse

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/articles/:article_id/like", h.createLikeHandler)
}

// createLikeHandler godoc
//
//	@Summary		Like an article
//	@Description	Like an article API.
//	@Tags			articles
//	@Produce		json
//	@Param			article_id	path	int	true	"Article ID"
// @Success		201
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		409			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/articles/{article_id}/like [post]
func (h *Handler) createLikeHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid article_id"})
		return
	}

	userID := auth.MustUserID(c)

	err = h.repo.Create(c.Request.Context(), articleID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errArticleNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		case errors.Is(err, errAlreadyLiked):
			c.JSON(http.StatusConflict, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to like article"})
		}
		return
	}

	c.Status(http.StatusCreated)
}
