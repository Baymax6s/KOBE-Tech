package list_replies

import "time"

type Reply struct {
	ID        int64
	ArticleID int64
	UserID    int64

	Content  string
	ParentID *int64

	Kind int

	IsBest bool

	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// kind定義（DB int前提）
const (
	KindComment  = 0
	KindQuestion = 1
	KindAnswer   = 2
)