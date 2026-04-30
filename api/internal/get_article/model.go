package article

import "time"

type Author struct {
	ID   int64
	Name string
}

type Article struct {
	ID        int64
	Title     string
	Content   string
	Author    Author
	CreatedAt time.Time
	UpdatedAt time.Time
}
