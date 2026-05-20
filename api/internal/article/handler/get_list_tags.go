package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/gin-gonic/gin"
)

type TagsErrorResponse struct {
	Message string `json:"message"`
} // @name server.tagsErrorResponse

type TagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.tagJSONResponse

type ListTagsJSONResponse struct {
	Tags []TagJSON `json:"tags"`
} // @name server.listTagsResponse

// listTagsHandler godoc
//
//	@Summary		List tags
//	@Description	Get all tag names. Used for tag candidates on the article list / create screens. No auth required.
//	@Tags			article
//	@Produce		json
//	@Success		200	{object}	ListTagsJSONResponse
//	@Failure		500	{object}	TagsErrorResponse
//	@Router			/api/tags [get]
func (h *Handler) listTagsHandler(c *gin.Context) {
	response, err := h.ListTags(c.Request.Context())
	if err != nil {
		log.Printf("list tags: %v", err)
		c.JSON(http.StatusInternalServerError, TagsErrorResponse{
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListTags(ctx context.Context) (ListTagsJSONResponse, error) {
	if h == nil || h.repo == nil {
		return ListTagsJSONResponse{}, errors.New("tags handler is not configured")
	}

	tags, err := h.repo.ListTags(ctx)
	if err != nil {
		return ListTagsJSONResponse{}, err
	}

	return newListTagsJSONResponse(tags), nil
}

func newListTagsJSONResponse(tags []article.Tag) ListTagsJSONResponse {
	response := ListTagsJSONResponse{
		Tags: make([]TagJSON, 0, len(tags)),
	}

	for _, tag := range tags {
		response.Tags = append(response.Tags, TagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}
