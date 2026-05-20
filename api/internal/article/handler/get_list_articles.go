package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type ArticleErrorResponse struct {
	Message string `json:"message"`
} // @name server.articleErrorResponse

type ArticleTagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleTagJSONResponse

type ArticleListItemJSON struct {
	ID         int64            `json:"id" binding:"required"`
	Title      string           `json:"title" binding:"required"`
	Content    string           `json:"content" binding:"required"`
	UserID     int64            `json:"user_id" binding:"required"`
	Tags       []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time        `json:"created_at" binding:"required"`
	UpdatedAt  time.Time        `json:"updated_at" binding:"required"`
	LikesCount int64            `json:"likes_count" binding:"required"`
	LikedByMe  bool             `json:"liked_by_me"`
} // @name server.articleJSONResponse

type ListArticlesJSONResponse struct {
	Articles []ArticleListItemJSON `json:"articles"`
} // @name server.listArticlesResponse

// listArticlesHandler godoc
//
//	@Summary		List articles
//	@Description	Get article list API. Supports filtering by tag name (AND).
//	@Tags			article
//	@Produce		json
//	@Param			tag	query		[]string	false	"Filter by tag name (multiple = AND)"	collectionFormat(multi)
//	@Success		200	{object}	ListArticlesJSONResponse
//	@Failure		400	{object}	ArticleErrorResponse
//	@Failure		500	{object}	ArticleErrorResponse
//	@Router			/api/articles [get]
func (h *Handler) listArticlesHandler(c *gin.Context) {
	userID, _ := auth.OptionalUserID(c)

	tagNames, err := normalizeFilterTagNames(c.QueryArray("tag"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ArticleErrorResponse{
			Message: "invalid tag query parameter",
		})
		return
	}

	response, err := h.ListArticles(c.Request.Context(), userID, tagNames)
	if err != nil {
		log.Printf("list articles: %v", err)
		c.JSON(http.StatusInternalServerError, ArticleErrorResponse{
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListArticles(ctx context.Context, userID int64, tagNames []string) (ListArticlesJSONResponse, error) {
	if h == nil || h.repo == nil {
		return ListArticlesJSONResponse{}, errors.New("article handler is not configured")
	}

	articles, err := h.repo.ListArticles(ctx, userID, tagNames)
	if err != nil {
		return ListArticlesJSONResponse{}, err
	}

	return newListArticlesJSONResponse(articles), nil
}

// normalizeFilterTagNames は ?tag=... のクエリパラメータを SQL の LOWER(t.name) 比較に揃える。
// trim / 長さ / 重複の検証は NormalizeTagNames に任せ、最後に lowercase だけ施す。
// 空文字は「絞り込み無し」の意味なので、エラーにせず除外してから渡す。
func normalizeFilterTagNames(raw []string) ([]string, error) {
	filtered := make([]string, 0, len(raw))
	for _, t := range raw {
		if strings.TrimSpace(t) == "" {
			continue
		}
		filtered = append(filtered, t)
	}
	if len(filtered) == 0 {
		return nil, nil
	}
	normalized, err := article.NormalizeTagNames(filtered)
	if err != nil {
		return nil, err
	}
	for i, n := range normalized {
		normalized[i] = strings.ToLower(n)
	}
	return normalized, nil
}

func newListArticlesJSONResponse(articles []article.Article) ListArticlesJSONResponse {
	response := ListArticlesJSONResponse{
		Articles: make([]ArticleListItemJSON, 0, len(articles)),
	}

	for _, item := range articles {
		response.Articles = append(response.Articles, ArticleListItemJSON{
			ID:         item.ID,
			Title:      item.Title,
			Content:    item.Content,
			UserID:     item.UserID,
			Tags:       newArticleTagJSONs(item.Tags),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			LikesCount: item.LikesCount,
			LikedByMe:  item.LikedByMe,
		})
	}

	return response
}

func newArticleTagJSONs(tags []article.Tag) []ArticleTagJSON {
	response := make([]ArticleTagJSON, 0, len(tags))
	for _, tag := range tags {
		response = append(response, ArticleTagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}
