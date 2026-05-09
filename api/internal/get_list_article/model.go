package article

import "time"

type Article struct {
	ID         int64
	Title      string
	Content    string
	UserID     int64
	Tags       []Tag
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LikesCount int64
}

type Tag struct {
	ID   int64
	Name string
}
