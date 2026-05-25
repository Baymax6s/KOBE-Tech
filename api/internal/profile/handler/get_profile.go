package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/minio"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	"github.com/gin-gonic/gin"
)

type ProfileJSON struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	ObjectKey string `json:"objectKey,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`

	ProfileCreatedAt *time.Time `json:"profile_created_at"`
	ProfileUpdatedAt *time.Time `json:"profile_updated_at"`
}

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.profileErrorResponse

// getProfileHandler
// @Summary      プロフィール取得
// @Description  ログインユーザーのプロフィール情報を取得する
// @Tags         profile
// @Accept       json
// @Produce      json
// @ID           profileList
// @Success      200  {object}  ProfileJSON
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/profile [get]
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

	profileJSON := newProfileJSON(user)

	if user.UserProfile.IsUploaded && user.UserProfile.ObjectKey.Valid {
		client, err := minio.NewClient()
		if err != nil {
			return ProfileJSON{}, fmt.Errorf("failed to connect to storage: %w", err)
		}

		bucketName := os.Getenv("MINIO_BUCKET_NAME")

		url, err := minio.GeneratePresignedGetURL(ctx, client, bucketName, profileJSON.ObjectKey)
		if err != nil {
			return ProfileJSON{}, fmt.Errorf("failed to generate presigned GET URL: %w", err)
		}

		profileJSON.AvatarURL = url
	}

	return profileJSON, nil
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

	var objectKey string
	if p.UserProfile.IsUploaded && p.UserProfile.ObjectKey.Valid {
		objectKey = p.UserProfile.ObjectKey.String
	}

	return ProfileJSON{
		ID:        p.User.ID,
		Name:      p.User.Name,
		Bio:       p.UserProfile.Bio.String,
		ObjectKey: objectKey,

		ProfileCreatedAt: profileCreatedAt,
		ProfileUpdatedAt: profileUpdatedAt,
	}
}
