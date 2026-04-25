package login

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
} // @name server.loginRequest

type LoginResponse struct {
	Token string `json:"token"`
} // @name server.loginResponse

type ErrorResponse struct {
	Message string `json:"message"`
} // @name server.loginErrorResponse

type Handler struct {
	repo   *Repository
	issuer *auth.Issuer
}

func NewHandler(repo *Repository, issuer *auth.Issuer) *Handler {
	return &Handler{repo: repo, issuer: issuer}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/auth/login", h.loginHandler)
}

// loginHandler godoc
//
//	@Summary		Login
//	@Description	Authenticates a user by name and password, and returns a JWT.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/auth/login [post]
func (h *Handler) loginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	token, err := h.Login(c.Request.Context(), req.Name, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidRequest):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		case errors.Is(err, errInvalidCredentials):
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to login"})
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
		return "", errors.New("login handler is not configured")
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
