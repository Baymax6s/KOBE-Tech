package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"
)

func (r *Repository) Create(ctx context.Context, articleID, userID int64) error {
	if r == nil || r.db == nil {
		return errors.New("like repository is not configured")
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
				if pqErr.Constraint == "fk_likes_article_id" {
					return ErrArticleNotFound
				}
				return ErrUserNotFound
			case "23505":
				return ErrAlreadyLiked
			}
		}
		return err
	}

	return nil
}
