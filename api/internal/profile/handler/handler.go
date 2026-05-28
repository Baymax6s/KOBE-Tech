package handler

import (
	articlerepository "github.com/Baymax6s/KOBE-Tech/api/internal/article/repository"
	"github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo        *repository.Repository
	articleRepo *articlerepository.Repository
}

func NewHandler(repo *repository.Repository, articleRepo *articlerepository.Repository) *Handler {
	return &Handler{repo: repo, articleRepo: articleRepo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) {
	router.GET("/profile/:user_id", h.getUserProfileHandler)
	router.GET("/profile/:user_id/articles", h.getUserArticlesHandler)
	router.GET("/profile/:user_id/liked-articles", h.getUserLikedArticlesHandler)
	router.GET("/profile/:user_id/best-answer-articles", h.getUserBestAnswerArticlesHandler)
	authRouter.PUT("/profile/bio", h.updateBioHandler)
	authRouter.POST("/profile/avatar/presign", h.presignAvatarHandler)
	authRouter.POST("/profile/avatar/complete", h.avatarUploadCompleteHandler)
}
