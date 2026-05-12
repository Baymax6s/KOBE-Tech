package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

// createArticleHandler godoc
//
//	@Summary		Create article
//	@Description	Create article API.
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createArticleRequest	true	"Create article request"
//	@Success		201		{object}	CreateArticleJSONResponse
//	@Failure		400		{object}	ArticleErrorResponse
//	@Failure		401		{object}	ArticleErrorResponse
//	@Failure		500		{object}	ArticleErrorResponse
//	@Security		BearerAuth
//	@Router			/api/articles [post]
func (h *Handler) createArticleHandler(c *gin.Context) {
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ArticleErrorResponse{Message: "invalid request body"})
		return
	}

	userID := auth.MustUserID(c)

	response, err := h.CreateArticle(c.Request.Context(), userID, req.Title, req.Content, req.TagNames)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidRequest), errors.Is(err, errInvalidTagName):
			c.JSON(http.StatusBadRequest, ArticleErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ArticleErrorResponse{Message: "failed to create article"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

var errInvalidRequest = errors.New("title and content are required")
var errInvalidTagName = errors.New("tags must contain tag names between 1 and 10 characters")

const maxTagNameLength = 10

func (h *Handler) CreateArticle(ctx context.Context, userID int64, title, content string, tagNames []string) (CreateArticleJSONResponse, error) {
	if h == nil || h.repo == nil {
		return CreateArticleJSONResponse{}, errors.New("post article handler is not configured")
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" || content == "" {
		return CreateArticleJSONResponse{}, errInvalidRequest
	}

	normalizedTagNames, err := normalizeTagNames(tagNames)
	if err != nil {
		return CreateArticleJSONResponse{}, err
	}

	item, err := h.repo.CreateArticle(ctx, title, content, userID, normalizedTagNames)
	if err != nil {
		return CreateArticleJSONResponse{}, err
	}

	return CreateArticleJSONResponse{
		ID:        item.ID,
		Title:     item.Title,
		Content:   item.Content,
		UserID:    item.UserID,
		Tags:      newArticleTagJSONs(item.Tags),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func normalizeTagNames(tagNames []string) ([]string, error) {
	normalizedTagNames := make([]string, 0, len(tagNames))
	seen := make(map[string]struct{}, len(tagNames))

	for _, tagName := range tagNames {
		normalizedTagName := strings.ToLower(strings.TrimSpace(tagName))
		if normalizedTagName == "" {
			return nil, errInvalidTagName
		}
		if utf8.RuneCountInString(normalizedTagName) > maxTagNameLength {
			return nil, errInvalidTagName
		}
		if _, ok := seen[normalizedTagName]; ok {
			continue
		}

		seen[normalizedTagName] = struct{}{}
		normalizedTagNames = append(normalizedTagNames, normalizedTagName)
	}

	return normalizedTagNames, nil
}
