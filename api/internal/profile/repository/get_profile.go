package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/profile"
)

func (r *Repository) FindByID(ctx context.Context, id int64) (profile.Profile, error) {
	if r == nil || r.db == nil {
		return profile.Profile{}, errors.New("repository not configured")
	}

	const query = `
        SELECT
        	users.id,
        	users.name,
        	COALESCE(user_profiles.bio, ''),
        	users.created_at,
        	users.updated_at
    	FROM users
    	LEFT JOIN user_profiles
        		ON user_profiles.user_id = users.id
    	WHERE users.id = $1
    `

	var user profile.Profile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.User.ID,
		&user.User.Name,
		&user.UserProfile.Bio,
		&user.User.CreatedAt,
		&user.User.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profile.Profile{}, ErrUserNotFound
		}
		return profile.Profile{}, err
	}

	return user, nil
}
