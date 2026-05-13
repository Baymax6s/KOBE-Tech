package handler

import (
	"github.com/Baymax6s/KOBE-Tech/api/internal/like/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/articles/:article_id/like", h.createLikeHandler)
	router.DELETE("/articles/:article_id/like", h.deleteLikeHandler)
}
