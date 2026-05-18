package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	"github.com/gin-gonic/gin"
)

type UpdateBioRequest struct {
	Bio string `json:"bio"`
} // @name server.updateBioRequest

// updateBioHandler godoc
//
// @Summary Update bio
// @Description ログインユーザーの自己紹介を更新する
// @Tags profile
// @Accept json
// @Produce json
// @Param request body UpdateBioRequest true "Update bio request"
// @Success 200 {object} ProfileJSON
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/profile/bio [put]
func (h *Handler) updateBioHandler(c *gin.Context) {
	userID := auth.MustUserID(c)

	var req UpdateBioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
		return
	}

	res, err := h.UpdateBio(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, profile.ErrInvalidBio),
			errors.Is(err, profile.ErrBioTooLong):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})

		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to update"})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateBio(ctx context.Context, userID int64, req UpdateBioRequest) (ProfileJSON, error) {
	if h == nil || h.repo == nil {
		return ProfileJSON{}, errors.New("handler not configured")
	}

	bio, err := profile.NormalizeBioInput(req.Bio)
	if err != nil {
		return ProfileJSON{}, err
	}

	if err := h.repo.UpdateBio(ctx, userID, bio); err != nil {
		return ProfileJSON{}, err
	}

	user, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return ProfileJSON{}, err
	}

	return newProfileJSON(user), nil
}
