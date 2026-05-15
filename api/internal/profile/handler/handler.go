package handler

import (
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) {
	authRouter.GET("/profile", h.getProfileHandler)
	authRouter.PUT("/profile", h.updateBioHandler)
}
