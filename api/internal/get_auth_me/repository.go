package me

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

func (r *Repository) FindByID(ctx context.Context, id int64) (User, error) {
	if r == nil || r.db == nil {
		return User{}, errors.New("me repository is not configured")
	}

	const query = `
		SELECT id, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
