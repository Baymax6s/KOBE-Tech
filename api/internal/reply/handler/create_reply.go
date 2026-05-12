package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
	"github.com/Baymax6s/KOBE-Tech/api/internal/reply/repository"
	"github.com/gin-gonic/gin"
)

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
			errors.Is(err, errBodyTooLong),
			errors.Is(err, errInvalidKind),
			errors.Is(err, repository.ErrInvalidParent),
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

const maxBodyLength = 2000

var (
	errInvalidBody = errors.New("body is required")
	errBodyTooLong = fmt.Errorf("body must be %d characters or less", maxBodyLength)
	errInvalidKind = errors.New("kind must be 'comment'")
)

func (h *Handler) CreateReply(ctx context.Context, articleID, userID int64, req createReplyRequest) (ReplyJSON, error) {
	if h == nil || h.repo == nil {
		return ReplyJSON{}, errors.New("reply handler is not configured")
	}

	body := strings.TrimSpace(req.Body)
	if body == "" {
		return ReplyJSON{}, errInvalidBody
	}
	if utf8.RuneCountInString(body) > maxBodyLength {
		return ReplyJSON{}, errBodyTooLong
	}

	// 今回スコープは comment のみ。kind 省略時は comment 扱い。
	kind := reply.KindComment
	if req.Kind != nil {
		parsed, ok := reply.ParseKind(*req.Kind)
		if !ok || parsed != reply.KindComment {
			return ReplyJSON{}, errInvalidKind
		}
		kind = parsed
	}

	item, err := h.repo.Create(ctx, articleID, userID, req.ParentID, kind, body)
	if err != nil {
		return ReplyJSON{}, err
	}

	return newReplyJSON(item), nil
}
