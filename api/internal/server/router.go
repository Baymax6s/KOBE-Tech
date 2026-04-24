package server

import (
	"database/sql"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	listarticle "github.com/Baymax6s/KOBE-Tech/api/internal/get_list_article"
	postarticle "github.com/Baymax6s/KOBE-Tech/api/internal/post_article"
	"github.com/gin-gonic/gin"
)

func NewHandler(db *sql.DB, validator *auth.Validator) http.Handler {
	listArticleHandler := listarticle.NewHandler(listarticle.NewRepository(db))
	postArticleHandler := postarticle.NewHandler(postarticle.NewRepository(db), validator)

	router := gin.Default()
	router.Use(corsMiddleware())
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	registerSwaggerRoutes(router)

	api := router.Group("/api")
	// Article
	listArticleHandler.RegisterRoutes(api)
	postArticleHandler.RegisterRoutes(api)

	return router
}
