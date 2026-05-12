package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
} // @name server.loginRequest

type LoginResponse struct {
	Token string `json:"token"`
} // @name server.loginResponse

type LoginErrorResponse struct {
	Message string `json:"message"`
} // @name server.loginErrorResponse

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
