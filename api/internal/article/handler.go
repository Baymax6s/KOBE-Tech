package article

import (
	"net/http"
	"strings"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type createArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Handler struct {
	repo      *Repository
	validator *auth.Validator
}

func NewHandler(repo *Repository, validator *auth.Validator) *Handler {
	return &Handler{repo: repo, validator: validator}
}

func (h *Handler) Create(c *gin.Context) {
	if h == nil || h.repo == nil || h.validator == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "article handler is not configured"})
		return
	}

	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "authorization header is required"})
		return
	}

	parts := strings.Fields(strings.TrimSpace(authorization))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
		return
	}
	userID, err := h.validator.ValidateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "title and content are required"})
		return
	}

	article, err := h.repo.Create(c.Request.Context(), req.Title, req.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}
