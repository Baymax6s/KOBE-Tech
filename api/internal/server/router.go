package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	article "github.com/Baymax6s/KOBE-Tech/api/internal/get_list_article"
)

func NewHandler(db *sql.DB) http.Handler {
	articleHandler := article.NewHandler(article.NewRepository(db))

	router := gin.Default()
	router.Use(corsMiddleware())
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	registerSwaggerRoutes(router)

	api := router.Group("/api")
	// article
	articleHandler.RegisterRoutes(api)

	return router
}
