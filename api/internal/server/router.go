package server

import (
	"database/sql"
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	me "github.com/Baymax6s/KOBE-Tech/api/internal/get_auth_me"
	listarticle "github.com/Baymax6s/KOBE-Tech/api/internal/get_list_article"
	postarticle "github.com/Baymax6s/KOBE-Tech/api/internal/post_article"
	login "github.com/Baymax6s/KOBE-Tech/api/internal/post_login"
	"github.com/gin-gonic/gin"
)

func NewHandler(db *sql.DB, validator *auth.Validator, issuer *auth.Issuer) http.Handler {
	listArticleHandler := listarticle.NewHandler(listarticle.NewRepository(db))
	postArticleHandler := postarticle.NewHandler(postarticle.NewRepository(db), validator)
	loginHandler := login.NewHandler(login.NewRepository(db), issuer)
	meHandler := me.NewHandler(me.NewRepository(db), validator)

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
	// Auth
	loginHandler.RegisterRoutes(api)
	meHandler.RegisterRoutes(api)

	return router
}
