package handler

import (
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth/repository"
	"github.com/gin-gonic/gin"
)

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
