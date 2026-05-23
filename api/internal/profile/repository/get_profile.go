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

    // 💡 user_profiles から object_key, is_uploaded, created_at, updated_at をそのまま取得します
    const query = `
        SELECT
            users.id,
            users.name,
            user_profiles.bio,
            user_profiles.object_key,
            user_profiles.is_uploaded,
            user_profiles.created_at,
            user_profiles.updated_at
        FROM users
        LEFT JOIN user_profiles
            ON user_profiles.user_id = users.id
        WHERE users.id = $1
    `

    var p profile.Profile

    // 💡 profile.go の sql.Null〜 型の定義と完全に一致させて Scan します
    // ※ is_uploaded だけは通常の bool なので、NULLの可能性がある LEFT JOIN 時の対策として
    // COALESCE(user_profiles.is_uploaded, false) にするか、Scan用に一度 NullBool を挟むのが安全です。
    // ここでは一番安全な COALESCE をSQL側に仕込んでスキャンします。
    const safeQuery = `
        SELECT
            users.id,
            users.name,
            user_profiles.bio,
            user_profiles.object_key,
            COALESCE(user_profiles.is_uploaded, false),
            user_profiles.created_at,
            user_profiles.updated_at
        FROM users
        LEFT JOIN user_profiles
            ON user_profiles.user_id = users.id
        WHERE users.id = $1
    `

    err := r.db.QueryRowContext(ctx, safeQuery, id).Scan(
        &p.User.ID,
        &p.User.Name,
        &p.UserProfile.Bio,
        &p.UserProfile.ObjectKey,
        &p.UserProfile.IsUploaded, // 💡 これで profile.go の bool 型に直接マッピングされます
        &p.UserProfile.CreatedAt,
        &p.UserProfile.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return profile.Profile{}, ErrUserNotFound
        }
        return profile.Profile{}, err
    }

    return p, nil
}