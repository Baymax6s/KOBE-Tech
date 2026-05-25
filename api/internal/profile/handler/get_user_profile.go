package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	"github.com/gin-gonic/gin"
)

type ProfileJSON struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Bio              string     `json:"bio"`
	ProfileCreatedAt *time.Time `json:"profile_created_at"`
	ProfileUpdatedAt *time.Time `json:"profile_updated_at"`
	IsOwner          bool       `json:"is_owner"`
} // @name server.profileJSON

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.profileErrorResponse

// getUserProfileHandler godoc
//
// @Summary Get user profile
// @Description ユーザーのプロフィールを取得する
// @Tags profile
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} ProfileJSON
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/profile/{user_id} [get]
func (h *Handler) getUserProfileHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid user_id"})
		return
	}

	currentUserID, _ := auth.OptionalUserID(c)

	res, err := h.GetProfile(c.Request.Context(), userID, currentUserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "failed to get profile",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProfile(ctx context.Context, targetUserID int64, currentUserID int64) (ProfileJSON, error) {
	p, err := h.repo.FindByID(ctx, targetUserID)
	if err != nil {
		return ProfileJSON{}, err
	}
	isOwner := currentUserID != 0 && currentUserID == targetUserID
	return newProfileJSON(p, isOwner), nil
}

func newProfileJSON(p profile.Profile, isOwner bool) ProfileJSON {
	var profileCreatedAt *time.Time
	if p.UserProfile.CreatedAt.Valid {
		t := p.UserProfile.CreatedAt.Time
		profileCreatedAt = &t
	}

	var profileUpdatedAt *time.Time
	if p.UserProfile.UpdatedAt.Valid {
		t := p.UserProfile.UpdatedAt.Time
		profileUpdatedAt = &t
	}

	return ProfileJSON{
		ID:               p.User.ID,
		Name:             p.User.Name,
		Bio:              p.UserProfile.Bio.String,
		ProfileCreatedAt: profileCreatedAt,
		ProfileUpdatedAt: profileUpdatedAt,
		IsOwner:          isOwner,
	}
}
