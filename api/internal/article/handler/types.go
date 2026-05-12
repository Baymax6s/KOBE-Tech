package handler

import "time"

type ArticleErrorResponse struct {
	Message string `json:"message"`
} // @name server.articleErrorResponse

type TagsErrorResponse struct {
	Message string `json:"message"`
} // @name server.tagsErrorResponse

type AuthorJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleAuthorJSONResponse

type ArticleTagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.articleTagJSONResponse

type ArticleListItemJSON struct {
	ID         int64            `json:"id" binding:"required"`
	Title      string           `json:"title" binding:"required"`
	Content    string           `json:"content" binding:"required"`
	UserID     int64            `json:"user_id" binding:"required"`
	Tags       []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time        `json:"created_at" binding:"required"`
	UpdatedAt  time.Time        `json:"updated_at" binding:"required"`
	LikesCount int64            `json:"likes_count" binding:"required"`
} // @name server.articleJSONResponse

type ListArticlesJSONResponse struct {
	Articles []ArticleListItemJSON `json:"articles"`
} // @name server.listArticlesResponse

type GetArticleJSONResponse struct {
	ID         int64            `json:"id" binding:"required"`
	Title      string           `json:"title" binding:"required"`
	Content    string           `json:"content" binding:"required"`
	Author     AuthorJSON       `json:"author" binding:"required"`
	Tags       []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt  time.Time        `json:"created_at" binding:"required"`
	UpdatedAt  time.Time        `json:"updated_at" binding:"required"`
	LikesCount int64            `json:"likes_count" binding:"required"`
} // @name server.getArticleJSONResponse

type createArticleRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	TagNames []string `json:"tags" minLength:"1" maxLength:"10"`
} // @name server.createArticleRequest

type CreateArticleJSONResponse struct {
	ID        int64            `json:"id" binding:"required"`
	Title     string           `json:"title" binding:"required"`
	Content   string           `json:"content" binding:"required"`
	UserID    int64            `json:"user_id" binding:"required"`
	Tags      []ArticleTagJSON `json:"tags" binding:"required"`
	CreatedAt time.Time        `json:"created_at" binding:"required"`
	UpdatedAt time.Time        `json:"updated_at" binding:"required"`
} // @name server.createArticleResponse

type TagJSON struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
} // @name server.tagJSONResponse

type ListTagsJSONResponse struct {
	Tags []TagJSON `json:"tags"`
} // @name server.listTagsResponse
