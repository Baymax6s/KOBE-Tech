package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/like/repository"
	"github.com/gin-gonic/gin"
)

type LikeCountResponse struct {
	LikesCount int64 `json:"likes_count"`
} // @name server.likeCountResponse

// deleteLikeHandler godoc
//
//	@Summary		Unlike an article
//	@Description	Unlike an article API.
//	@Tags			articles
//	@Produce		json
//	@Param			article_id	path	int	true	"Article ID"
//
//	@Success		200			{object}	LikeCountResponse
//
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/articles/{article_id}/like [delete]
func (h *Handler) deleteLikeHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid article_id"})
		return
	}

	userID := auth.MustUserID(c)

	err = h.repo.Delete(c.Request.Context(), articleID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotLiked):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to unlike article"})
		}
		return
	}

	count, err := h.repo.CountByArticle(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to count likes"})
		return
	}

	c.JSON(http.StatusOK, LikeCountResponse{LikesCount: count})
}
