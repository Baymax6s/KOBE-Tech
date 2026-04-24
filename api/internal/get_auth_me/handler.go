package me

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type MeResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name server.meResponse

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.meErrorResponse

type Handler struct {
	repo      *Repository
	validator *auth.Validator
}

func NewHandler(repo *Repository, validator *auth.Validator) *Handler {
	return &Handler{repo: repo, validator: validator}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/auth/me", h.meHandler)
}

// meHandler godoc
//
//	@Summary		Get current user
//	@Description	Returns the authenticated user identified by the JWT in the Authorization header.
//	@Tags			auth
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Success		200				{object}	MeResponse
//	@Failure		401				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/api/auth/me [get]
func (h *Handler) meHandler(c *gin.Context) {
	res, err := h.Me(c.Request.Context(), c.GetHeader("Authorization"))
	if err != nil {
		switch {
		case errors.Is(err, errInvalidAuthorization),
			errors.Is(err, errInvalidToken),
			errors.Is(err, errUserNotFound):
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get current user"})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

var (
	errInvalidAuthorization = errors.New("invalid authorization header")
	errInvalidToken         = errors.New("invalid token")
	errUserNotFound         = errors.New("user not found")
)

func (h *Handler) Me(ctx context.Context, authorization string) (MeResponse, error) {
	if h == nil || h.repo == nil || h.validator == nil {
		return MeResponse{}, errors.New("me handler is not configured")
	}

	parts := strings.Fields(strings.TrimSpace(authorization))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return MeResponse{}, errInvalidAuthorization
	}

	userID, err := h.validator.ValidateToken(parts[1])
	if err != nil {
		return MeResponse{}, errInvalidToken
	}

	user, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MeResponse{}, errUserNotFound
		}
		return MeResponse{}, err
	}

	return MeResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
