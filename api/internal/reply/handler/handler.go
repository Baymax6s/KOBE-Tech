package handler

import (
	"github.com/Baymax6s/KOBE-Tech/api/internal/reply/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) {
	router.GET("/articles/:article_id/replies", h.listRepliesHandler)
	authRouter.POST("/articles/:article_id/replies", h.createReplyHandler)
	authRouter.POST("/replies/:reply_id/best", h.setBestAnswerHandler)
	authRouter.PUT("/replies/:reply_id/best", h.unsetBestAnswerHandler)
}
