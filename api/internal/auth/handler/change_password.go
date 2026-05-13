package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
} // @name server.changePasswordRequest

type ChangePasswordResponse struct {
	Message string `json:"message"`
} // @name server.changePasswordResponse

type ChangePasswordErrorResponse struct {
	Message string `json:"message"`
} // @name server.changePasswordErrorResponse

// changePasswordHandler godoc
//
//	@Summary		Change password
//	@Description	Changes the authenticated user's password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ChangePasswordRequest	true	"Change password request"
//	@Success		200		{object}	ChangePasswordResponse
//	@Failure		400		{object}	ChangePasswordErrorResponse
//	@Failure		401		{object}	ChangePasswordErrorResponse
//	@Failure		500		{object}	ChangePasswordErrorResponse
//	@Security		BearerAuth
//	@Router			/api/auth/password [put]
func (h *Handler) changePasswordHandler(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ChangePasswordErrorResponse{Message: "invalid request body"})
		return
	}

	userID := auth.MustUserID(c)

	err := h.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, errPasswordValidation):
			c.JSON(http.StatusBadRequest, ChangePasswordErrorResponse{Message: err.Error()})
		case errors.Is(err, errUserNotFound), errors.Is(err, errCurrentPasswordIncorrect):
			c.JSON(http.StatusUnauthorized, ChangePasswordErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ChangePasswordErrorResponse{Message: "failed to change password"})
		}
		return
	}

	c.JSON(http.StatusOK, ChangePasswordResponse{Message: "パスワードを変更しました"})
}

var errPasswordValidation = errors.New("validation error")
var errCurrentPasswordIncorrect = errors.New("current password is incorrect")

func (h *Handler) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	if h == nil || h.repo == nil {
		return errors.New("auth handler is not configured")
	}

	if currentPassword == "" || newPassword == "" {
		return errors.Join(errPasswordValidation, errors.New("current_password and new_password are required"))
	}
	if len(newPassword) < 8 {
		return errors.Join(errPasswordValidation, errors.New("new_password must be at least 8 characters"))
	}

	user, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errUserNotFound
		}
		return err
	}

	if err := auth.ComparePassword(user.PasswordHash, currentPassword); err != nil {
		return errCurrentPasswordIncorrect
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return h.repo.UpdatePassword(ctx, userID, newHash)
}
