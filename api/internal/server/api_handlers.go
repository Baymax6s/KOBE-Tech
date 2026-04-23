package server

import (
	"net/http"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
)

type listArticlesResponse = article.ListArticlesJSONResponse
type articleErrorResponse = article.ErrorJSONResponse

type notImplementedResponse struct {
	Message  string `json:"message" example:"feature is not implemented yet"`
	NextStep string `json:"next_step" example:"internal/{domain}/{handler,service,repository}.go"`
}

type loginRequest struct {
	Email    string `json:"email" format:"email" example:"user@example.com"`
	Password string `json:"password" format:"password" example:"change-me"`
}

type createArticleRequest struct {
	Title string `json:"title" example:"First article"`
	Body  string `json:"body" example:"Article body"`
}

// loginHandler godoc
//
//	@Summary		Login
//	@Description	Login API.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest			true	"Login request"
//	@Failure		501		{object}	notImplementedResponse
//	@Router			/api/auth/login [post]
func (s *apiServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "login", "internal/auth/{handler,service,repository}.go")
}

// listArticlesHandler godoc
//
//	@Summary		List articles
//	@Description	Get article list API.
//	@Tags			articles
//	@Produce		json
//	@Success		200	{object}	listArticlesResponse
//	@Failure		500	{object}	articleErrorResponse
//	@Router			/api/articles [get]
func (s *apiServer) listArticlesHandler(w http.ResponseWriter, r *http.Request) {
	s.articleHandler.ListArticles(w, r)
}

// createArticleHandler godoc
//
//	@Summary		Create article
//	@Description	Create article API.
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createArticleRequest	true	"Create article request"
//	@Failure		501		{object}	notImplementedResponse
//	@Router			/api/articles [post]
func (s *apiServer) createArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "create article", "internal/article/{handler,service,repository}.go")
}
