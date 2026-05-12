package article

import "time"

type Author struct {
	ID   int64
	Name string
}

type Article struct {
	ID         int64
	Title      string
	Content    string
	UserID     int64
	Author     Author
	Tags       []Tag
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LikesCount int64
}
