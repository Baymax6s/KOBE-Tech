package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type ProfileErrorResponse struct {
	Message string `json:"message"`
} // @name server.profileErrorResponse

type ProfileArticleTagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.profileArticleTagJSON

type ProfileArticleListItemJSON struct {
	ID         int64                   `json:"id" binding:"required"`
	Title      string                  `json:"title" binding:"required"`
	Content    string                  `json:"content" binding:"required"`
	UserID     int64                   `json:"user_id" binding:"required"`
	Tags       []ProfileArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time               `json:"created_at" binding:"required"`
	UpdatedAt  time.Time               `json:"updated_at" binding:"required"`
	LikesCount int64                   `json:"likes_count" binding:"required"`
	LikedByMe  bool                    `json:"liked_by_me"`
} // @name server.profileArticleListItemJSON

type ProfileListArticlesJSONResponse struct {
	Articles []ProfileArticleListItemJSON `json:"articles"`
} // @name server.profileListArticlesJSONResponse

// listUserArticlesHandler godoc
//
//	@Summary		List articles by user
//	@Description	指定ユーザーの記事一覧を取得する
//	@Tags			profile
//	@Produce		json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	ProfileListArticlesJSONResponse
//	@Failure		400		{object}	ProfileErrorResponse
//	@Failure		500		{object}	ProfileErrorResponse
//	@Router			/api/profile/{user_id}/articles [get]
func (h *Handler) listUserArticlesHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, ProfileErrorResponse{Message: "invalid user_id"})
		return
	}

	currentUserID, _ := auth.OptionalUserID(c)

	articles, err := h.repo.ListArticlesByUser(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		log.Printf("list user articles: %v", err)
		c.JSON(http.StatusInternalServerError, ProfileErrorResponse{Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toProfileListArticlesJSONResponse(articles))
}

func toProfileListArticlesJSONResponse(articles []article.Article) ProfileListArticlesJSONResponse {
	resp := ProfileListArticlesJSONResponse{
		Articles: make([]ProfileArticleListItemJSON, 0, len(articles)),
	}
	for _, item := range articles {
		resp.Articles = append(resp.Articles, ProfileArticleListItemJSON{
			ID:         item.ID,
			Title:      item.Title,
			Content:    item.Content,
			UserID:     item.UserID,
			Tags:       toProfileArticleTagJSONs(item.Tags),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			LikesCount: item.LikesCount,
			LikedByMe:  item.LikedByMe,
		})
	}
	return resp
}

func toProfileArticleTagJSONs(tags []article.Tag) []ProfileArticleTagJSON {
	result := make([]ProfileArticleTagJSON, 0, len(tags))
	for _, tag := range tags {
		result = append(result, ProfileArticleTagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return result
}
