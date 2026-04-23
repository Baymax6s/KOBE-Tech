package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
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

	c.Header("Content-Type", contentType)
	http.ServeContent(c.Writer, c.Request, name, time.Time{}, bytes.NewReader(body))
}
