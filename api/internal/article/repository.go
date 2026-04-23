package article

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Repository) Create(ctx context.Context, title, content string, userID int64) (*Article, error) {
	const query = `INSERT INTO articles (title, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, title, content, user_id, created_at, updated_at`

	article := &Article{}
	err := r.db.QueryRowContext(ctx, query, title, content, userID).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return article, nil
}
