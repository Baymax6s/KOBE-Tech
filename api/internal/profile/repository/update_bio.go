package repository

import (
	"context"
	"errors"
)

func (r *Repository) UpdateBio(ctx context.Context, id int64, bio string) error {
	if r == nil || r.db == nil {
		return errors.New("repository not configured")
	}

	const query = `
        UPDATE users
        SET bio = $1, updated_at = NOW()
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, bio, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
