package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/auth"
)

func (r *Repository) FindByID(ctx context.Context, id int64) (auth.User, error) {
	if r == nil || r.db == nil {
		return auth.User{}, errors.New("auth repository is not configured")
	}

	const query = `
        SELECT id, name, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user auth.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return auth.User{}, err
	}

	return user, nil
}
