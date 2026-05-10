package reply

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

type ListRepliesJSONResponse struct {
	Replies []ReplyJSON `json:"replies" binding:"required"`
} // @name server.listRepliesResponse

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
	router.GET("/articles/:article_id/replies", h.listRepliesHandler)
}

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
		return ListRepliesJSONResponse{}, errors.New("get list reply handler is not configured")
	}

	replies, err := h.repo.ListByArticleID(ctx, articleID)
	if err != nil {
		return ListRepliesJSONResponse{}, err
	}

	items := make([]ReplyJSON, 0, len(replies))
	for _, reply := range replies {
		items = append(items, ReplyJSON{
			ID:        reply.ID,
			ArticleID: reply.ArticleID,
			ParentID:  reply.ParentID,
			Kind:      reply.Kind.String(),
			Body:      reply.Body,
			UserID:    reply.UserID,
			UserName:  reply.UserName,
			CreatedAt: reply.CreatedAt,
			UpdatedAt: reply.UpdatedAt,
		})
	}

	return ListRepliesJSONResponse{Replies: items}, nil
}
