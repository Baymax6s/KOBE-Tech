package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
