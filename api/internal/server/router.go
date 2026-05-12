package server

import (
	"database/sql"
	"net/http"

	articlehandler "github.com/Baymax6s/KOBE-Tech/api/internal/article/handler"
	articlerepository "github.com/Baymax6s/KOBE-Tech/api/internal/article/repository"
	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
	authhandler "github.com/Baymax6s/KOBE-Tech/api/internal/auth/handler"
	authrepository "github.com/Baymax6s/KOBE-Tech/api/internal/auth/repository"
	likehandler "github.com/Baymax6s/KOBE-Tech/api/internal/like/handler"
	likerepository "github.com/Baymax6s/KOBE-Tech/api/internal/like/repository"
	profilehandler "github.com/Baymax6s/KOBE-Tech/api/internal/profile/handler"
	profilerepository "github.com/Baymax6s/KOBE-Tech/api/internal/profile/repository"
	replyhandler "github.com/Baymax6s/KOBE-Tech/api/internal/reply/handler"
	replyrepository "github.com/Baymax6s/KOBE-Tech/api/internal/reply/repository"
	"github.com/gin-gonic/gin"
)

func NewHandler(db *sql.DB, validator *auth.Validator, issuer *auth.Issuer) http.Handler {
	articleHandler := articlehandler.NewHandler(articlerepository.NewRepository(db))
	authHandler := authhandler.NewHandler(authrepository.NewRepository(db), issuer)
	likeHandler := likehandler.NewHandler(likerepository.NewRepository(db))
	replyHandler := replyhandler.NewHandler(replyrepository.NewRepository(db))
	profileHandler := profilehandler.NewHandler(profilerepository.NewRepository(db))

	router := gin.Default()
	router.Use(corsMiddleware())
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	registerSwaggerRoutes(router)

	api := router.Group("/api")

	authRequired := api.Group("", auth.RequireUser(validator))
	articleHandler.RegisterRoutes(api, authRequired)
	authHandler.RegisterRoutes(api, authRequired)
	likeHandler.RegisterRoutes(authRequired)
	replyHandler.RegisterRoutes(api, authRequired)
	profileHandler.RegisterRoutes(api, authRequired)

	return router
}
