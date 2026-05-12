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
}

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

func (h *Handler) UpdateBio(ctx context.Context, userID int64, req UpdateBioRequest) (gin.H, error) {
    if h == nil || h.repo == nil {
        return nil, errors.New("handler not configured")
    }

    bio, err := profile.NormalizeBioInput(req.Bio)
    if err != nil {
        return nil, err
    }

    if err := h.repo.UpdateBio(ctx, userID, bio); err != nil {
        return nil, err
    }

    return gin.H{"message": "updated"}, nil
}