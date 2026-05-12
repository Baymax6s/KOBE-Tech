package handler

import "time"

type createReplyRequest struct {
	ParentID *int64  `json:"parent_id,omitempty"`
	Kind     *string `json:"kind,omitempty"`
	Body     string  `json:"body" binding:"required"`
} // @name server.createReplyRequest

type ReplyJSON struct {
	ID        int64     `json:"id" binding:"required"`
	ArticleID int64     `json:"article_id" binding:"required"`
	ParentID  *int64    `json:"parent_id,omitempty"`
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
