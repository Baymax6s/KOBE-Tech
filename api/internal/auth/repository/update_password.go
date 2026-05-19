package repository

import (
	"context"
	"errors"
)

func (r *Repository) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	if r == nil || r.db == nil {
		return errors.New("auth repository is not configured")
	}

	const query = `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}
