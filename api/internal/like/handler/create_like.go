package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/like/repository"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.likeErrorResponse

type LikeResponse struct {
	LikesCount int64 `json:"likes_count"`
	LikedByMe  bool  `json:"liked_by_me"`
} // @name server.likeResponse

// createLikeHandler godoc
//
//	@Summary		Like an article
//	@Description	Like an article API.
//	@Tags			like
//	@Produce		json
//	@Param			article_id	path	int	true	"Article ID"
//
//	@Success		200			{object}	LikeResponse
//
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
		case errors.Is(err, repository.ErrArticleNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		case errors.Is(err, repository.ErrAlreadyLiked):
			c.JSON(http.StatusConflict, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to like article"})
		}
		return
	}

	count, err := h.repo.CountByArticle(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to count likes"})
		return
	}

	c.JSON(http.StatusOK, LikeResponse{LikesCount: count, LikedByMe: true})
}
