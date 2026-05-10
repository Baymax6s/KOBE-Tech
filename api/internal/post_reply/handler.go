package reply

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type createReplyRequest struct {
	ParentID *int64  `json:"parent_id"`
	Kind     *string `json:"kind"`
	Body     string  `json:"body"`
} // @name server.createReplyRequest

type ReplyJSON struct {
	ID        int64     `json:"id" binding:"required"`
	ArticleID int64     `json:"article_id" binding:"required"`
	ParentID  *int64    `json:"parent_id"`
	Kind      string    `json:"kind" binding:"required" enums:"comment,question,answer"`
	Body      string    `json:"body" binding:"required"`
	UserID    int64     `json:"user_id" binding:"required"`
	UserName  string    `json:"user_name" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
} // @name server.replyJSONResponse

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.replyErrorResponse

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/articles/:article_id/replies", h.createReplyHandler)
}

// createReplyHandler godoc
//
//	@Summary		Create a reply (comment) on an article
//	@Description	記事 / 既存コメントへのコメントを投稿する。今回のスコープは kind = comment のみ。
//	@Tags			replies
//	@Accept			json
//	@Produce		json
//	@Param			article_id	path		int					true	"Article ID"
//	@Param			request		body		createReplyRequest	true	"Create reply request"
//	@Success		201			{object}	ReplyJSON
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/articles/{article_id}/replies [post]
func (h *Handler) createReplyHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid article_id"})
		return
	}

	var req createReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	userID := auth.MustUserID(c)

	response, err := h.CreateReply(c.Request.Context(), articleID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidBody),
			errors.Is(err, errInvalidKind),
			errors.Is(err, errInvalidParent),
			errors.Is(err, errParentMismatch):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		case errors.Is(err, errArticleNotFound), errors.Is(err, errParentNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to create reply"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

var (
	errInvalidBody = errors.New("body is required")
	errInvalidKind = errors.New("kind must be 'comment'")
)

const maxBodyLength = 2000

func (h *Handler) CreateReply(ctx context.Context, articleID, userID int64, req createReplyRequest) (ReplyJSON, error) {
	if h == nil || h.repo == nil {
		return ReplyJSON{}, errors.New("post reply handler is not configured")
	}

	body := strings.TrimSpace(req.Body)
	if body == "" || len(body) > maxBodyLength {
		return ReplyJSON{}, errInvalidBody
	}

	// 今回スコープは comment のみ。kind 省略時は comment 扱い。
	kind := KindComment
	if req.Kind != nil {
		parsed, ok := ParseKind(*req.Kind)
		if !ok || parsed != KindComment {
			return ReplyJSON{}, errInvalidKind
		}
		kind = parsed
	}

	reply, err := h.repo.Create(ctx, articleID, userID, req.ParentID, kind, body)
	if err != nil {
		return ReplyJSON{}, err
	}

	return ReplyJSON{
		ID:        reply.ID,
		ArticleID: reply.ArticleID,
		ParentID:  reply.ParentID,
		Kind:      reply.Kind.String(),
		Body:      reply.Body,
		UserID:    reply.UserID,
		UserName:  reply.UserName,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}, nil
}
