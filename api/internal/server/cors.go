package server

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	allowedMethods        = "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD"
	defaultHeaders        = "Accept, Authorization, Content-Type, Origin"
	defaultAllowedOrigins = "https://baymax6s.github.io,https://vue-cjne.onrender.com"
	defaultMaxAgeSeconds  = 0
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
		headers.Set("Access-Control-Max-Age", corsMaxAge())

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
			headers.Set("Cache-Control", "no-store")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string) bool {
	if isConfiguredAllowedOrigin(origin) {
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

func isConfiguredAllowedOrigin(origin string) bool {
	for _, allowedOrigin := range configuredAllowedOrigins() {
		if strings.EqualFold(allowedOrigin, origin) {
			return true
		}
	}

	return false
}

func configuredAllowedOrigins() []string {
	origins := []string{defaultAllowedOrigins}
	if extraOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); extraOrigins != "" {
		origins = append(origins, extraOrigins)
	}

	allowedOrigins := make([]string, 0, len(origins))
	for _, group := range origins {
		for _, origin := range strings.Split(group, ",") {
			origin = strings.TrimSpace(origin)
			if origin == "" {
				continue
			}
			allowedOrigins = append(allowedOrigins, origin)
		}
	}

	return allowedOrigins
}

func corsMaxAge() string {
	if value := strings.TrimSpace(os.Getenv("CORS_MAX_AGE_SECONDS")); value != "" {
		seconds, err := strconv.Atoi(value)
		if err == nil && seconds >= 0 {
			return strconv.Itoa(seconds)
		}
	}

	return strconv.Itoa(defaultMaxAgeSeconds)
}
