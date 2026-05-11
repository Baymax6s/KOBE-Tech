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
	Author     Author
	Tags       []Tag
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LikesCount int64
}

type Tag struct {
	ID   int64
	Name string
}
