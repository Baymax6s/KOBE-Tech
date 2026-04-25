package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const userIDContextKey = "auth.userID"

type ErrorResponse struct {
	Message string `json:"message"`
}

func RequireUser(v *Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := extractBearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "invalid authorization header",
			})
			return
		}

		userID, err := v.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "invalid token",
			})
			return
		}

		c.Set(userIDContextKey, userID)
		c.Next()
	}
}

// MustUserID は RequireUser 適用下でのみ呼ぶこと。未適用なら panic する。
func MustUserID(c *gin.Context) int64 {
	v, exists := c.Get(userIDContextKey)
	if !exists {
		panic("auth.MustUserID called without auth.RequireUser middleware")
	}
	userID, ok := v.(int64)
	if !ok {
		panic("auth.MustUserID: userID in context is not int64")
	}
	return userID
}

func extractBearerToken(header string) (string, bool) {
	parts := strings.Fields(strings.TrimSpace(header))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}
