package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	"github.com/Baymax6s/KOBE-Tech/api/internal/minio"
	"github.com/gin-gonic/gin"

	minioSDK "github.com/minio/minio-go/v7"
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

// presignAvatarHandler godoc
//
// @Summary Generate profile avatar upload URL
// @Description ログインユーザーのアバターアップロード用のPresigned PUT URLを発行する
// @Tags profile
// @Accept json
// @Produce json
// @Param request body PresignAvatarRequest true "Presign avatar request"
// @Success 200 {object} PresignAvatarResponse
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to update profile"})
		return
	}

	client, err := minio.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to connect to storage"})
		return
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	url, err := minio.GeneratePresignedPutURL(c.Request.Context(), client, bucketName, objectKey)
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
// @Success 200 {object} AvatarUploadCompleteResponse
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

	userIDStr := strconv.FormatInt(userID, 10)
	expectedJPG := fmt.Sprintf("avatar/%s.jpg", userIDStr)
	expectedPNG := fmt.Sprintf("avatar/%s.png", userIDStr)

	if req.ObjectKey != expectedJPG && req.ObjectKey != expectedPNG {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid objectKey for this user"})
		return
	}

	client, err := minio.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to connect to storage"})
		return
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	objInfo, err := client.StatObject(c.Request.Context(), bucketName, req.ObjectKey, minioSDK.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "uploaded file not found or inaccessible"})
		return
	}

	if objInfo.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "uploaded file exceeds 2MB limit"})
		return
	}

	if err := h.repo.UpsertUserProfile(c.Request.Context(), userID, req.ObjectKey, true); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to mark upload complete"})
		return
	}

	c.JSON(http.StatusOK, AvatarUploadCompleteResponse{ObjectKey: req.ObjectKey, Uploaded: true})
}
