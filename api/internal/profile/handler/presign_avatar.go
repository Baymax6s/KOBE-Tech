package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/minio"
	"github.com/gin-gonic/gin"
)

type PresignAvatarRequest struct {
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
}

type PresignAvatarResponse struct {
	URL       string `json:"url"`
	ObjectKey string `json:"objectKey"`
}

type PresignAvatarCompleteRequest struct {
	ObjectKey string `json:"objectKey"`
}

type AvatarUploadCompleteResponse struct {
	ObjectKey string `json:"objectKey"`
	Uploaded  bool   `json:"uploaded"`
}

type PresignGetAvatarRequest struct {
	ObjectKey string `json:"objectKey"`
}

type PresignGetAvatarResponse struct {
	URL string `json:"url"`
}

// presignAvatarHandler godoc
//
// @Summary Generate profile avatar upload URL
// @Description ログインユーザーのアバターアップロード用のPresigned PUT URLを発行する
// @Tags profile
// @Accept json
// @Produce json
// @Param request body PresignAvatarRequest true "Presign avatar request"
// @Success 200 {object} handler.PresignAvatarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/profile/avatar/presign [post]
func (h *Handler) presignAvatarHandler(c *gin.Context) {
	var req PresignAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
		return
	}

	if req.Size <= 0 || req.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "image size must be 2MB or less"})
		return
	}

	if req.ContentType != "" && req.ContentType != "image/png" && req.ContentType != "image/jpeg" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "unsupported content type"})
		return
	}

	userID := auth.MustUserID(c)

	ext := "png"
	if req.ContentType == "image/jpeg" {
		ext = "jpg"
	}

	objectKey := fmt.Sprintf("avatar/%s.%s", strconv.FormatInt(userID, 10), ext)

	if err := h.repo.UpsertUserProfile(c.Request.Context(), userID, objectKey, false); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: fmt.Sprintf("upsert failed: %v", err)})
		return
	}

	client, err := minio.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to connect to storage"})
		return
	}

	url, err := minio.GeneratePresignedPutURL(client, objectKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to generate presigned URL"})
		return
	}

	c.JSON(http.StatusOK, PresignAvatarResponse{URL: url, ObjectKey: objectKey})
}

// avatarUploadCompleteHandler godoc
//
// @Summary Mark profile avatar upload complete
// @Description フロントエンドからの完了通知を受けてDBのアップロード状態を true にする
// @Tags profile
// @Accept json
// @Produce json
// @Param request body PresignAvatarCompleteRequest true "Upload complete request"
// @Success 200 {object} handler.AvatarUploadCompleteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/profile/avatar/complete [post]
func (h *Handler) avatarUploadCompleteHandler(c *gin.Context) {
	var req PresignAvatarCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
		return
	}

	if req.ObjectKey == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "objectKey is required"})
		return
	}

	userID := auth.MustUserID(c)
	if err := h.repo.UpsertUserProfile(c.Request.Context(), userID, req.ObjectKey, true); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to mark upload complete"})
		return
	}

	c.JSON(http.StatusOK, AvatarUploadCompleteResponse{ObjectKey: req.ObjectKey, Uploaded: true})
}

// presignGetAvatarHandler godoc
//
// @Summary Generate profile avatar download URL
// @Description ログインユーザーのアバター閲覧用のPresigned GET URLを発行する
// @Tags profile
// @Accept json
// @Produce json
// @Param request body handler.PresignGetAvatarRequest true "Presign GET avatar request"
// @Success 200 {object} handler.PresignGetAvatarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/profile/avatar/download [post]
func (h *Handler) presignGetAvatarHandler(c *gin.Context) {
	var req PresignGetAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
		return
	}

	if req.ObjectKey == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "objectKey is required"})
		return
	}

	_ = auth.MustUserID(c)

	client, err := minio.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to connect to storage"})
		return
	}

	url, err := minio.GeneratePresignedGetURL(client, req.ObjectKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to generate presigned GET URL"})
		return
	}

	c.JSON(http.StatusOK, PresignGetAvatarResponse{URL: url})
}