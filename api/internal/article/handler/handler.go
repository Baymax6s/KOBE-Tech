package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/Baymax6s/KOBE-Tech/api/internal/article/repository"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type ArticleErrorResponse struct {
	Message string `json:"message"`
} // @name server.articleErrorResponse

type TagsErrorResponse struct {
	Message string `json:"message"`
} // @name server.tagsErrorResponse

type AuthorJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleAuthorJSONResponse

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
} // @name server.articleJSONResponse

type ListArticlesJSONResponse struct {
	Articles []ArticleListItemJSON `json:"articles"`
} // @name server.listArticlesResponse

type GetArticleJSONResponse struct {
	ID         int64            `json:"id" binding:"required"`
	Title      string           `json:"title" binding:"required"`
	Content    string           `json:"content" binding:"required"`
	Author     AuthorJSON       `json:"author" binding:"required"`
	Tags       []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time        `json:"created_at" binding:"required"`
	UpdatedAt  time.Time        `json:"updated_at" binding:"required"`
	LikesCount int64            `json:"likes_count" binding:"required"`
} // @name server.getArticleJSONResponse

type createArticleRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	TagNames []string `json:"tags" minLength:"1" maxLength:"10"`
} // @name server.createArticleRequest

type CreateArticleJSONResponse struct {
	ID        int64            `json:"id" binding:"required"`
	Title     string           `json:"title" binding:"required"`
	Content   string           `json:"content" binding:"required"`
	UserID    int64            `json:"user_id" binding:"required"`
	Tags      []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt time.Time        `json:"created_at" binding:"required"`
	UpdatedAt time.Time        `json:"updated_at" binding:"required"`
} // @name server.createArticleResponse

type TagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.tagJSONResponse

type ListTagsJSONResponse struct {
	Tags []TagJSON `json:"tags"`
} // @name server.listTagsResponse

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) {
	router.GET("/articles", h.listArticlesHandler)
	router.GET("/articles/:article_id", h.getArticleHandler)
	authRouter.POST("/articles", h.createArticleHandler)
	authRouter.GET("/tags", h.listTagsHandler)
}

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

// getArticleHandler godoc
//
//	@Summary		Get article
//	@Description	Get article detail API.
//	@Tags			articles
//	@Produce		json
//	@Param			article_id	path		int	true	"Article ID"
//	@Success		200			{object}	GetArticleJSONResponse
//	@Failure		400			{object}	ArticleErrorResponse
//	@Failure		404			{object}	ArticleErrorResponse
//	@Failure		500			{object}	ArticleErrorResponse
//	@Router			/api/articles/{article_id} [get]
func (h *Handler) getArticleHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ArticleErrorResponse{Message: "invalid article_id"})
		return
	}

	response, err := h.GetArticle(c.Request.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, errArticleNotFound):
			c.JSON(http.StatusNotFound, ArticleErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ArticleErrorResponse{Message: "failed to get article"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

var errArticleNotFound = errors.New("article not found")

func (h *Handler) GetArticle(ctx context.Context, articleID int64) (GetArticleJSONResponse, error) {
	if h == nil || h.repo == nil {
		return GetArticleJSONResponse{}, errors.New("get article handler is not configured")
	}

	item, err := h.repo.FindArticleByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetArticleJSONResponse{}, errArticleNotFound
		}
		return GetArticleJSONResponse{}, err
	}

	return GetArticleJSONResponse{
		ID:      item.ID,
		Title:   item.Title,
		Content: item.Content,
		Author: AuthorJSON{
			ID:   item.Author.ID,
			Name: item.Author.Name,
		},
		Tags:       newArticleTagJSONs(item.Tags),
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
		LikesCount: item.LikesCount,
	}, nil
}

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
