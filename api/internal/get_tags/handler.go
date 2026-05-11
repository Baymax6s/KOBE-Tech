package tags

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.tagsErrorResponse

type TagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.tagJSONResponse

type ListTagsJSONResponse struct {
	Tags []TagJSON `json:"tags"`
} // @name server.listTagsResponse

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/tags", h.listTagsHandler)
}

// listTagsHandler godoc
//
//	@Summary		List tags
//	@Description	Get tag candidates for creating articles.
//	@Tags			tags
//	@Produce		json
//	@Success		200	{object}	ListTagsJSONResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/tags [get]
func (h *Handler) listTagsHandler(c *gin.Context) {
	response, err := h.ListTags(c.Request.Context())
	if err != nil {
		log.Printf("list tags: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
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

	tags, err := h.repo.List(ctx)
	if err != nil {
		return ListTagsJSONResponse{}, err
	}

	return newListTagsJSONResponse(tags), nil
}

func newListTagsJSONResponse(tags []Tag) ListTagsJSONResponse {
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
