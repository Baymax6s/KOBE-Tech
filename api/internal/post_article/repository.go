package article

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, title, content string, userID int64) (Article, error) {
	if r == nil || r.db == nil {
		return Article{}, errors.New("post article repository is not configured")
	}

	const query = `
		INSERT INTO articles (title, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, title, content, user_id, created_at, updated_at
	`

	var article Article
	err := r.db.QueryRowContext(ctx, query, title, content, userID).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return Article{}, err
	}

	return article, nil
}
