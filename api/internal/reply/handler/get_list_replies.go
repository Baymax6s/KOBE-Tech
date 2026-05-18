package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
	"github.com/gin-gonic/gin"
)

type ReplyJSON struct {
	ID        int64     `json:"id" binding:"required"`
	ArticleID int64     `json:"article_id" binding:"required"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	Kind      string    `json:"kind" binding:"required" enums:"comment,question,answer"`
	Body      string    `json:"body" binding:"required"`
	IsBest    bool      `json:"is_best" binding:"required"`
	UserID    int64     `json:"user_id" binding:"required"`
	UserName  string    `json:"user_name" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
} // @name server.replyJSONResponse

type ListRepliesJSONResponse struct {
	Replies []ReplyJSON `json:"replies" binding:"required"`
} // @name server.listRepliesResponse

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.replyErrorResponse

// listRepliesHandler godoc
//
//	@Summary		List replies of an article
//	@Description	記事に紐づく返信（コメント / 質問 / 回答）を全件取得する。
//	@Tags			replies
//	@Produce		json
//	@Param			article_id	path		int	true	"Article ID"
//	@Success		200			{object}	ListRepliesJSONResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/articles/{article_id}/replies [get]
func (h *Handler) listRepliesHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid article_id"})
		return
	}

	response, err := h.ListReplies(c.Request.Context(), articleID)
	if err != nil {
		log.Printf("list replies: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to list replies"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListReplies(ctx context.Context, articleID int64) (ListRepliesJSONResponse, error) {
	if h == nil || h.repo == nil {
		return ListRepliesJSONResponse{}, errors.New("reply handler is not configured")
	}

	replies, err := h.repo.ListByArticleID(ctx, articleID)
	if err != nil {
		return ListRepliesJSONResponse{}, err
	}

	return ListRepliesJSONResponse{Replies: newReplyJSONs(replies)}, nil
}

func newReplyJSONs(replies []reply.Reply) []ReplyJSON {
	items := make([]ReplyJSON, 0, len(replies))
	for _, item := range replies {
		items = append(items, newReplyJSON(item))
	}
	return items
}

func newReplyJSON(item reply.Reply) ReplyJSON {
	return ReplyJSON{
		ID:        item.ID,
		ArticleID: item.ArticleID,
		ParentID:  item.ParentID,
		Kind:      string(item.Kind),
		Body:      item.Body,
		IsBest:    item.IsBest,
		UserID:    item.UserID,
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
