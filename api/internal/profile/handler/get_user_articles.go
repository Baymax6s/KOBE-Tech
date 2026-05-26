package handler

import (
	"log"
	"net/http"
	"strconv"

	articlehandler "github.com/Baymax6s/KOBE-Tech/api/internal/article/handler"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

// getUserArticlesHandler godoc
//
//	@Summary		List articles posted by a user
//	@Description	指定ユーザーが投稿した記事一覧を新しい順に取得する
//	@Tags			profile
//	@Produce		json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	handler.ListArticlesJSONResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/profile/{user_id}/articles [get]
func (h *Handler) getUserArticlesHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid user_id"})
		return
	}

	viewerID, _ := auth.OptionalUserID(c)

	articles, err := h.articleRepo.ListArticlesByAuthor(c.Request.Context(), userID, viewerID)
	if err != nil {
		log.Printf("list user articles: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to list articles"})
		return
	}

	c.JSON(http.StatusOK, articlehandler.NewListArticlesJSONResponse(articles))
}

// getUserLikedArticlesHandler godoc
//
//	@Summary		List articles liked by a user
//	@Description	指定ユーザーがいいねした記事一覧をいいねした順に取得する
//	@Tags			profile
//	@Produce		json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	handler.ListArticlesJSONResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/profile/{user_id}/liked-articles [get]
func (h *Handler) getUserLikedArticlesHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid user_id"})
		return
	}

	viewerID, _ := auth.OptionalUserID(c)

	articles, err := h.articleRepo.ListLikedArticles(c.Request.Context(), userID, viewerID)
	if err != nil {
		log.Printf("list user liked articles: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to list liked articles"})
		return
	}

	c.JSON(http.StatusOK, articlehandler.NewListArticlesJSONResponse(articles))
}
