package server

import "net/http"

type notImplementedResponse struct {
	Message  string `json:"message" example:"login is not implemented yet"`
	NextStep string `json:"next_step" example:"internal/auth/{handler,service,repository}.go"`
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
//	@Router			/api/v1/auth/login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "login", "internal/auth/{handler,service,repository}.go")
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
//	@Router			/api/v1/articles [post]
func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "create article", "internal/article/{handler,service,repository}.go")
}
