package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/gin-gonic/gin"
)

// タグ内ランキングで返す上位件数。TOP3 固定。
const popularArticlesPerTag = 3

type PopularArticleJSON struct {
	ArticleID  int64  `json:"article_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	LikesCount int64  `json:"likes_count" binding:"required"`
	Rank       int64  `json:"rank" binding:"required"`
} // @name server.popularArticleJSONResponse

type TagRankingJSON struct {
	TagID    int64                `json:"tag_id" binding:"required"`
	TagName  string               `json:"tag_name" binding:"required"`
	Articles []PopularArticleJSON `json:"articles" binding:"required"`
} // @name server.tagRankingJSONResponse

type PopularArticlesByTagJSONResponse struct {
	Rankings []TagRankingJSON `json:"rankings"`
} // @name server.popularArticlesByTagResponse

// listPopularArticlesByTagHandler godoc
//
//	@Summary		List popular articles per tag
//	@Description	Get top liked articles within each tag (category). No auth required.
//	@Tags			article
//	@Produce		json
//	@Success		200	{object}	PopularArticlesByTagJSONResponse
//	@Failure		500	{object}	ArticleErrorResponse
//	@Router			/api/tags/popular-articles [get]
func (h *Handler) listPopularArticlesByTagHandler(c *gin.Context) {
	response, err := h.PopularArticlesByTag(c.Request.Context())
	if err != nil {
		log.Printf("popular articles by tag: %v", err)
		c.JSON(http.StatusInternalServerError, ArticleErrorResponse{
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) PopularArticlesByTag(ctx context.Context) (PopularArticlesByTagJSONResponse, error) {
	if h == nil || h.repo == nil {
		return PopularArticlesByTagJSONResponse{}, errors.New("article handler is not configured")
	}

	rankings, err := h.repo.PopularArticlesByTag(ctx, popularArticlesPerTag)
	if err != nil {
		return PopularArticlesByTagJSONResponse{}, err
	}

	return newPopularArticlesByTagJSONResponse(rankings), nil
}

func newPopularArticlesByTagJSONResponse(rankings []article.TagRanking) PopularArticlesByTagJSONResponse {
	response := PopularArticlesByTagJSONResponse{
		Rankings: make([]TagRankingJSON, 0, len(rankings)),
	}

	for _, ranking := range rankings {
		articles := make([]PopularArticleJSON, 0, len(ranking.Articles))
		for _, item := range ranking.Articles {
			articles = append(articles, PopularArticleJSON{
				ArticleID:  item.ArticleID,
				Title:      item.Title,
				LikesCount: item.LikesCount,
				Rank:       item.Rank,
			})
		}

		response.Rankings = append(response.Rankings, TagRankingJSON{
			TagID:    ranking.TagID,
			TagName:  ranking.TagName,
			Articles: articles,
		})
	}

	return response
}
