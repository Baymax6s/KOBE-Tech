package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// listArticlesHandler godoc
//
//	@Summary		List articles
//	@Description	Get article list API.
//	@Tags			articles
//	@Produce		json
//	@Success		200	{object}	ListArticlesJSONResponse
//	@Failure		500	{object}	ArticleErrorResponse
//	@Router			/api/articles [get]
func (h *Handler) listArticlesHandler(c *gin.Context) {
	response, err := h.ListArticles(c.Request.Context())
	if err != nil {
		log.Printf("list articles: %v", err)
		c.JSON(http.StatusInternalServerError, ArticleErrorResponse{
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListArticles(ctx context.Context) (ListArticlesJSONResponse, error) {
	if h == nil || h.repo == nil {
		return ListArticlesJSONResponse{}, errors.New("article handler is not configured")
	}

	articles, err := h.repo.ListArticles(ctx)
	if err != nil {
		return ListArticlesJSONResponse{}, err
	}

	return newListArticlesJSONResponse(articles), nil
}
