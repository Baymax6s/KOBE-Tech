package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"
)

func (r *Repository) UpdateBio(ctx context.Context, id int64, bio string) error {
	if r == nil || r.db == nil {
		return errors.New("repository not configured")
	}

	// プロフィール行は移行時に bio が NULL だったユーザーには作られていない。
	// 「行が無い（まだ書いていない）」と「bio が空文字」を区別し、行が無ければ作成・あれば更新する。
	const query = `
		INSERT INTO user_profiles (user_id, bio, is_uploaded, created_at, updated_at)
		VALUES ($2, $1, FALSE, NOW(), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET bio = EXCLUDED.bio, updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, bio, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" && pqErr.Constraint == "fk_user_profiles_user_id" {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}
