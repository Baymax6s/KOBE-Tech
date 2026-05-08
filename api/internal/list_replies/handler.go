package list_replies

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/articles/:article_id/replies", h.getRepliesHandler)
}

func (h *Handler) getRepliesHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid article_id"})
		return
	}

	replyList, err := h.repo.FindByArticleID(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get replies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"replies": replyList,
	})
}