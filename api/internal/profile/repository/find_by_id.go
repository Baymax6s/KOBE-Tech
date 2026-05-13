package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/profile"
)

func (r *Repository) FindByID(ctx context.Context, id int64) (profile.User, error) {
    if r == nil || r.db == nil {
        return profile.User{}, errors.New("repository not configured")
    }

    const query = `
        SELECT id, name, COALESCE(bio, ''), created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var user profile.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Bio,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return profile.User{}, ErrUserNotFound
    }
    return profile.User{}, err
}


    return user, nil
}