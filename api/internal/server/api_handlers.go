package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type articleErrorResponse = messageResponse

type articleJSONResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type listArticlesResponse struct {
	Articles []articleJSONResponse `json:"articles"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type notImplementedResponse struct {
	Message  string `json:"message"`
	NextStep string `json:"next_step"`
}

type loginRequest struct {
	Email    string `json:"email" format:"email"`
	Password string `json:"password" format:"password"`
}

type createArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
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
func (s *apiServer) loginHandler(c *gin.Context) {
	writeNotImplemented(c, "login", "internal/auth/{handler,service,repository}.go")
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
func (s *apiServer) listArticlesHandler(c *gin.Context) {
	response, err := s.articleHandler.ListArticles(c.Request.Context())
	if err != nil {
		log.Printf("list articles: %v", err)
		writeJSON(c, http.StatusInternalServerError, messageResponse{
			Message: "internal server error",
		})
		return
	}

	writeJSON(c, http.StatusOK, response)
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
func (s *apiServer) createArticleHandler(c *gin.Context) {
	writeNotImplemented(c, "create article", "internal/article/{handler,service,repository}.go")
}
