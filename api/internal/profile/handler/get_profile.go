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
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// getProfileHandler godoc
//
// @Summary Get profile
// @Description ログインユーザーのプロフィール取得
// @Tags profile
// @Produce json
// @Success 200 {object} ProfileJSON
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
			Message: err.Error(),
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

func newProfileJSON(u profile.User) ProfileJSON {
	return ProfileJSON{
		ID:        u.ID,
		Name:      u.Name,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
