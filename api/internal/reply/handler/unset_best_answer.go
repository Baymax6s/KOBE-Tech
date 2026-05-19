package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/reply/repository"
	"github.com/gin-gonic/gin"
)

type UnsetBestAnswerResponse struct {
	ReplyID int64 `json:"reply_id"`
	IsBest  bool  `json:"is_best"`
} // @name server.unsetBestAnswerResponse

// unsetBestAnswerHandler godoc
//
//	@Summary		Unmark a reply as best answer
//	@Description	ベストアンサー指定を解除する。質問者のみ可能。
//	@Tags			replies
//	@Produce		json
//	@Param			reply_id	path		int	true	"Reply ID"
//	@Success		200			{object}	UnsetBestAnswerResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		409			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/replies/{reply_id}/best [put]
func (h *Handler) unsetBestAnswerHandler(c *gin.Context) {
	replyID, err := strconv.ParseInt(c.Param("reply_id"), 10, 64)
	if err != nil || replyID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid reply_id"})
		return
	}

	userID := auth.MustUserID(c)

	err = h.repo.UnsetBestAnswer(c.Request.Context(), replyID, userID)
	if err != nil {
		code := http.StatusInternalServerError
		switch {
		case errors.Is(err, repository.ErrReplyNotFound):
			code = http.StatusNotFound
		case errors.Is(err, repository.ErrNotAnswer):
			code = http.StatusBadRequest
		case errors.Is(err, repository.ErrNotBestAnswer):
			code = http.StatusConflict
		case errors.Is(err, repository.ErrParentNotFound):
			code = http.StatusNotFound
		case errors.Is(err, repository.ErrNotQuestionAuthor):
			code = http.StatusForbidden
		}
		c.JSON(code, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UnsetBestAnswerResponse{ReplyID: replyID, IsBest: false})
}
