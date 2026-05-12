package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth/repository"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
} // @name server.loginRequest

type LoginResponse struct {
	Token string `json:"token"`
} // @name server.loginResponse

type MeResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name server.meResponse

type LoginErrorResponse struct {
	Message string `json:"message"`
} // @name server.loginErrorResponse

type MeErrorResponse struct {
	Message string `json:"message"`
} // @name server.meErrorResponse

type Handler struct {
	repo   *repository.Repository
	issuer *auth.Issuer
}

func NewHandler(repo *repository.Repository, issuer *auth.Issuer) *Handler {
	return &Handler{repo: repo, issuer: issuer}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) {
	router.POST("/auth/login", h.loginHandler)
	authRouter.GET("/auth/me", h.meHandler)
}

// loginHandler godoc
//
//	@Summary		Login
//	@Description	Authenticates a user by name and password, and returns a JWT.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Login request"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	LoginErrorResponse
//	@Failure		401		{object}	LoginErrorResponse
//	@Failure		500		{object}	LoginErrorResponse
//	@Router			/api/auth/login [post]
func (h *Handler) loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginErrorResponse{Message: "invalid request body"})
		return
	}

	token, err := h.Login(c.Request.Context(), req.Name, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidRequest):
			c.JSON(http.StatusBadRequest, LoginErrorResponse{Message: err.Error()})
		case errors.Is(err, errInvalidCredentials):
			c.JSON(http.StatusUnauthorized, LoginErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, LoginErrorResponse{Message: "failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

var (
	errInvalidRequest     = errors.New("name and password are required")
	errInvalidCredentials = errors.New("invalid name or password")
	errUserNotFound       = errors.New("user not found")
)

func (h *Handler) Login(ctx context.Context, name, password string) (string, error) {
	if h == nil || h.repo == nil || h.issuer == nil {
		return "", errors.New("auth handler is not configured")
	}

	name = strings.TrimSpace(name)
	if name == "" || password == "" {
		return "", errInvalidRequest
	}

	user, err := h.repo.FindByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errInvalidCredentials
		}
		return "", err
	}

	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		return "", errInvalidCredentials
	}

	return h.issuer.Issue(user.ID)
}

// meHandler godoc
//
//	@Summary		Get current user
//	@Description	Returns the authenticated user identified by the JWT in the Authorization header.
//	@Tags			auth
//	@Produce		json
//	@Success		200				{object}	MeResponse
//	@Failure		401				{object}	MeErrorResponse
//	@Failure		500				{object}	MeErrorResponse
//	@Security		BearerAuth
//	@Router			/api/auth/me [get]
func (h *Handler) meHandler(c *gin.Context) {
	userID := auth.MustUserID(c)

	res, err := h.Me(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, errUserNotFound):
			c.JSON(http.StatusUnauthorized, MeErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, MeErrorResponse{Message: "failed to get current user"})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Me(ctx context.Context, userID int64) (MeResponse, error) {
	if h == nil || h.repo == nil {
		return MeResponse{}, errors.New("auth handler is not configured")
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
