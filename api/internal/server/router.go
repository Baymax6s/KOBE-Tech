package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	article "github.com/Baymax6s/KOBE-Tech/api/internal/get_list_article"
)

type apiServer struct {
	articleHandler *article.Handler
}

func NewHandler(db *sql.DB) http.Handler {
	server := &apiServer{
		articleHandler: article.NewHandler(article.NewRepository(db)),
	}

	router := gin.Default()

	registerSwaggerRoutes(router)

	api := router.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/login", server.loginHandler)
	api.GET("/articles", server.listArticlesHandler)
	api.POST("/articles", server.createArticleHandler)

	return router
}

func writeNotImplemented(c *gin.Context, feature, nextStep string) {
	writeJSON(c, http.StatusNotImplemented, notImplementedResponse{
		Message:  feature + " is not implemented yet",
		NextStep: nextStep,
	})
}

func writeJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}
