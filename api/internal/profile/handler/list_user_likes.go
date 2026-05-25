package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

// listUserLikesHandler godoc
//
//	@Summary		List liked articles by user
//	@Description	指定ユーザーがいいねした記事一覧を取得する
//	@Tags			profile
//	@Produce		json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	ProfileListArticlesJSONResponse
//	@Failure		400		{object}	ProfileErrorResponse
//	@Failure		500		{object}	ProfileErrorResponse
//	@Router			/api/profile/{user_id}/likes [get]
func (h *Handler) listUserLikesHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, ProfileErrorResponse{Message: "invalid user_id"})
		return
	}

	currentUserID, _ := auth.OptionalUserID(c)

	articles, err := h.repo.ListLikedArticles(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		log.Printf("list user likes: %v", err)
		c.JSON(http.StatusInternalServerError, ProfileErrorResponse{Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toProfileListArticlesJSONResponse(articles))
}
