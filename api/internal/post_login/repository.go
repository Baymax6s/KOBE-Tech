package postlogin

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

func (r *Repository) FindByName(ctx context.Context, name string) (*User, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("auth repository is not configured")
	}
	query := `
		SELECT id, name, password_hash, created_at, updated_at
		FROM users
		WHERE name = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&user.ID,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
