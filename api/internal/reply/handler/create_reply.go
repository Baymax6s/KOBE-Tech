package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
	"github.com/Baymax6s/KOBE-Tech/api/internal/reply/repository"
	"github.com/gin-gonic/gin"
)

type createReplyRequest struct {
	ParentID *int64  `json:"parent_id,omitempty"`
	Kind     *string `json:"kind,omitempty"`
	Body     string  `json:"body" binding:"required"`
} // @name server.createReplyRequest

// createReplyHandler godoc
//
//	@Summary		Create a reply on an article
//	@Description	記事 / 既存返信への返信を投稿する。kind は comment / question / answer のいずれか。ルート投稿は comment か question、コメント配下は comment、質問・回答配下は answer のみ受け付ける。
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
		case errors.Is(err, reply.ErrInvalidBody),
			errors.Is(err, reply.ErrBodyTooLong),
			errors.Is(err, reply.ErrInvalidKind),
			errors.Is(err, repository.ErrInvalidParent),
			errors.Is(err, repository.ErrInvalidRootKind),
			errors.Is(err, repository.ErrParentMismatch):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		case errors.Is(err, repository.ErrArticleNotFound), errors.Is(err, repository.ErrParentNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to create reply"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) CreateReply(ctx context.Context, articleID, userID int64, req createReplyRequest) (ReplyJSON, error) {
	if h == nil || h.repo == nil {
		return ReplyJSON{}, errors.New("reply handler is not configured")
	}

	body, kind, err := reply.NormalizeCreateInput(req.Body, req.Kind)
	if err != nil {
		return ReplyJSON{}, err
	}
	item, err := h.repo.Create(ctx, articleID, userID, req.ParentID, kind, body)
	if err != nil {
		return ReplyJSON{}, err
	}

	return newReplyJSON(item), nil
}
