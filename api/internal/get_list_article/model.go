package article

import "time"

type Article struct {
	ID        int64
	Title     string
	Content   string
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
