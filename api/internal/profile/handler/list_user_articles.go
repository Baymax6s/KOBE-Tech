package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	articlehandler "github.com/Baymax6s/KOBE-Tech/api/internal/article/handler"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

// listUserArticlesHandler godoc
//
//	@Summary		List articles by user
//	@Description	指定ユーザーの記事一覧を取得する
//	@Tags			profile
//	@Produce		json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	articlehandler.ListArticlesJSONResponse
//	@Failure		400		{object}	articlehandler.ArticleErrorResponse
//	@Failure		500		{object}	articlehandler.ArticleErrorResponse
//	@Router			/api/profile/{user_id}/articles [get]
func (h *Handler) listUserArticlesHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, articlehandler.ArticleErrorResponse{Message: "invalid user_id"})
		return
	}

	currentUserID, _ := auth.OptionalUserID(c)

	articles, err := h.repo.ListArticlesByUser(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		log.Printf("list user articles: %v", err)
		c.JSON(http.StatusInternalServerError, articlehandler.ArticleErrorResponse{Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toListArticlesJSONResponse(articles))
}

func toListArticlesJSONResponse(articles []article.Article) articlehandler.ListArticlesJSONResponse {
	resp := articlehandler.ListArticlesJSONResponse{
		Articles: make([]articlehandler.ArticleListItemJSON, 0, len(articles)),
	}
	for _, item := range articles {
		resp.Articles = append(resp.Articles, articlehandler.ArticleListItemJSON{
			ID:         item.ID,
			Title:      item.Title,
			Content:    item.Content,
			UserID:     item.UserID,
			Tags:       toArticleTagJSONs(item.Tags),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			LikesCount: item.LikesCount,
			LikedByMe:  item.LikedByMe,
		})
	}
	return resp
}

func toArticleTagJSONs(tags []article.Tag) []articlehandler.ArticleTagJSON {
	result := make([]articlehandler.ArticleTagJSON, 0, len(tags))
	for _, tag := range tags {
		result = append(result, articlehandler.ArticleTagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return result
}
