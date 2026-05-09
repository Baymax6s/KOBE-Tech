package list_replies

import "time"

type Reply struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	UserID    int64     `json:"user_id"`

	Content  string `json:"content"`
	ParentID *int64 `json:"parent_id"`

	Kind int `json:"kind"`

	IsBest bool `json:"is_best"`

	UserName string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// kind定義（DB int前提）
const (
	KindComment  = 0
	KindQuestion = 1
	KindAnswer   = 2
)