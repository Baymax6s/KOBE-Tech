package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD"
	defaultHeaders = "Accept, Authorization, Content-Type, Origin"
	maxAgeSeconds  = "600"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if !isAllowedOrigin(origin) {
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.Next()
			return
		}

		headers := c.Writer.Header()
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Methods", allowedMethods)
		headers.Set("Access-Control-Max-Age", maxAgeSeconds)

		allowHeaders := c.GetHeader("Access-Control-Request-Headers")
		if allowHeaders == "" {
			allowHeaders = defaultHeaders
		}
		headers.Set("Access-Control-Allow-Headers", allowHeaders)

		if strings.EqualFold(c.GetHeader("Access-Control-Request-Private-Network"), "true") {
			headers.Add("Vary", "Access-Control-Request-Private-Network")
			headers.Set("Access-Control-Allow-Private-Network", "true")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string) bool {
	if origin == "https://baymax6s.github.io" {
		return true
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	switch strings.ToLower(parsed.Hostname()) {
	case "localhost", "127.0.0.1":
		return true
	default:
		return false
	}
}
