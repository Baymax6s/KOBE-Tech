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
	postArticleHandler := postarticle.NewHandler(postarticle.NewRepository(db))
	loginHandler := login.NewHandler(login.NewRepository(db), issuer)
	meHandler := me.NewHandler(me.NewRepository(db))

	router := gin.Default()
	router.Use(corsMiddleware())
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	registerSwaggerRoutes(router)

	api := router.Group("/api")
	listArticleHandler.RegisterRoutes(api)
	loginHandler.RegisterRoutes(api)

	authRequired := api.Group("", auth.RequireUser(validator))
	postArticleHandler.RegisterRoutes(authRequired)
	meHandler.RegisterRoutes(authRequired)

	return router
}
