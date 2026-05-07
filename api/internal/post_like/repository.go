package like

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	errArticleNotFound = errors.New("article not found")
	errAlreadyLiked    = errors.New("already liked")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, articleID, userID int64) error {
	if r == nil || r.db == nil {
		return errors.New("post like repository is not configured")
	}

	const query = `
		INSERT INTO likes (article_id, user_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query, articleID, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23503":
				return errArticleNotFound
			case "23505":
				return errAlreadyLiked
			}
		}
		return err
	}

	return nil
}
