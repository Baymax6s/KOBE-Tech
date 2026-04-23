package server

import (
	"database/sql"
	"net/http"

	createarticle "github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	listarticle "github.com/Baymax6s/KOBE-Tech/api/internal/get_list_article"
	"github.com/gin-gonic/gin"
)

func NewHandler(db *sql.DB, validator *auth.Validator) http.Handler {
	listArticleHandler := listarticle.NewHandler(listarticle.NewRepository(db))
	createArticleHandler := createarticle.NewHandler(createarticle.NewRepository(db), validator)

	router := gin.Default()
	router.Use(corsMiddleware())
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	registerSwaggerRoutes(router)

	api := router.Group("/api")
	listArticleHandler.RegisterRoutes(api)
	api.POST("/articles", createArticleHandler.Create)

	return router
}
