package post_reply

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/articles/:article_id/replies", h.postReplyHandler)
}

type request struct {
	Content  string `json:"content"`
	ParentID *int64 `json:"parent_id"`
	Kind     int    `json:"kind"`
}

func (h *Handler) postReplyHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid article_id",
		})
		return
	}

	var req request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid body",
		})
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "content is required",
		})
		return
	}

	userID := auth.MustUserID(c)

	var parent *Reply

	if req.ParentID != nil {
		p, err := h.repo.FindByID(c.Request.Context(), *req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "parent not found",
			})
			return
		}

		parent = &p
	}

	if err := validateKind(parent, req.Kind); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	reply, err := h.repo.Create(
		c.Request.Context(),
		articleID,
		userID,
		req.Content,
		req.ParentID,
		req.Kind,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create reply",
		})
		return
	}

	c.JSON(http.StatusCreated, reply)
}

func validateKind(parent *Reply, kind int) error {
	// 記事直下
	if parent == nil {
		if kind != KindComment && kind != KindQuestion {
			return errors.New("root reply must be comment or question")
		}

		return nil
	}

	switch parent.Kind {

	case KindComment:
		if kind != KindComment {
			return errors.New("comment can only have comment replies")
		}

	case KindQuestion:
		if kind != KindAnswer {
			return errors.New("question can only have answer replies")
		}

	case KindAnswer:
		if kind != KindAnswer {
			return errors.New("answer can only have answer replies")
		}
	}

	return nil
}