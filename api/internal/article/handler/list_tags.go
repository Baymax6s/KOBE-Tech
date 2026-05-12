package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// listTagsHandler godoc
//
//	@Summary		List tags
//	@Description	Get tag candidates for creating articles.
//	@Tags			tags
//	@Produce		json
//	@Success		200	{object}	ListTagsJSONResponse
//	@Failure		401	{object}	TagsErrorResponse
//	@Failure		500	{object}	TagsErrorResponse
//	@Security		BearerAuth
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
