package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	"github.com/gin-gonic/gin"
)

type ProfileJSON struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`

	UserCreatedAt    time.Time `json:"user_created_at"`
	UserUpdatedAt    time.Time `json:"user_updated_at"`
	ProfileCreatedAt *time.Time `json:"profile_created_at"`
	ProfileUpdatedAt *time.Time `json:"profile_updated_at"`
}

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.profileErrorResponse

// getProfileHandler godoc
//
// @Summary Get profile
// @Description ログインユーザーのプロフィール取得
// @Tags profile
// @Produce json
// @Success 200 {object} profile_handler.ProfileJSON
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/profile [get]
func (h *Handler) getProfileHandler(c *gin.Context) {
	userID := auth.MustUserID(c)

	res, err := h.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "failed to get profile",
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProfile(ctx context.Context, userID int64) (ProfileJSON, error) {
	user, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return ProfileJSON{}, err
	}
	return newProfileJSON(user), nil
}

func newProfileJSON(p profile.Profile) ProfileJSON {

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
		ID:   p.User.ID,
		Name: p.User.Name,
		Bio:  p.UserProfile.Bio.String,

		UserCreatedAt: p.User.CreatedAt,
		UserUpdatedAt: p.User.UpdatedAt,

		ProfileCreatedAt: profileCreatedAt,
		ProfileUpdatedAt: profileUpdatedAt,
	}
}