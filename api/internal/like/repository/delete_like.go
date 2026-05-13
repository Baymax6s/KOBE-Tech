package repository

import (
	"context"
	"errors"
)

func (r *Repository) Delete(ctx context.Context, articleID, userID int64) error {
	if r == nil || r.db == nil {
		return errors.New("like repository is not configured")
	}

	const query = `
		DELETE FROM likes
		WHERE article_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, articleID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotLiked
	}

	return nil
}

func (r *Repository) CountByArticle(ctx context.Context, articleID int64) (int64, error) {
	if r == nil || r.db == nil {
		return 0, errors.New("like repository is not configured")
	}

	const query = `
		SELECT COUNT(*) FROM likes WHERE article_id = $1
	`

	var count int64
	err := r.db.QueryRowContext(ctx, query, articleID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
