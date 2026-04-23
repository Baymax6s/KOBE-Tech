package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	swaggerassets "github.com/Baymax6s/KOBE-Tech/api/swagger"
)

func registerSwaggerRoutes(router gin.IRouter) {
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
	router.GET("/swagger/", func(c *gin.Context) {
		serveSwaggerFile(c, "index.html", "text/html; charset=utf-8")
	})
	router.GET("/swagger/index.html", func(c *gin.Context) {
		serveSwaggerFile(c, "index.html", "text/html; charset=utf-8")
	})
	router.GET("/swagger/openapi.yml", func(c *gin.Context) {
		serveSwaggerFile(c, "openapi.yml", "application/yaml")
	})
	router.GET("/swagger/swagger.json", func(c *gin.Context) {
		serveSwaggerFile(c, "swagger.json", "application/json; charset=utf-8")
	})
}

func serveSwaggerFile(c *gin.Context, name, contentType string) {
	body, err := fs.ReadFile(swaggerassets.Files, name)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("swagger asset %q could not be loaded", name))
		return
	}

	if name == "swagger.json" {
		body, err = withSwaggerRequestTarget(c.Request, body)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("swagger asset %q could not be prepared", name))
			return
		}
	}

	c.Header("Content-Type", contentType)
	http.ServeContent(c.Writer, c.Request, name, time.Time{}, bytes.NewReader(body))
}

func withSwaggerRequestTarget(req *http.Request, body []byte) ([]byte, error) {
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}

	doc["host"] = firstHeaderValue(req.Header.Get("X-Forwarded-Host"))
	if doc["host"] == "" {
		doc["host"] = req.Host
	}
	doc["schemes"] = []string{swaggerScheme(req)}

	return json.MarshalIndent(doc, "", "    ")
}

func swaggerScheme(req *http.Request) string {
	if scheme := firstHeaderValue(req.Header.Get("X-Forwarded-Proto")); scheme != "" {
		return scheme
	}
	if req.TLS != nil {
		return "https"
	}

	return "http"
}

func firstHeaderValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	head, _, _ := strings.Cut(value, ",")
	return strings.TrimSpace(head)
}
