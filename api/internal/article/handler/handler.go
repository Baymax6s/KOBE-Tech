package handler

import (
	"github.com/Baymax6s/KOBE-Tech/api/internal/article/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) {
	router.GET("/articles", h.listArticlesHandler)
	router.GET("/articles/:article_id", h.getArticleHandler)
	router.GET("/tags", h.listTagsHandler)
	authRouter.POST("/articles", h.createArticleHandler)
}
