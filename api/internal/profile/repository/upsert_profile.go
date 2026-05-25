package repository

import (
    "context"
    "fmt"
)

func (r *Repository) UpsertUserProfile(ctx context.Context, userID int64, objectKey string, isUploaded bool) error {
    const query = `
        INSERT INTO user_profiles (user_id, object_key, is_uploaded, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        ON CONFLICT (user_id)
        DO UPDATE SET
            object_key = EXCLUDED.object_key,
            is_uploaded = EXCLUDED.is_uploaded,
            updated_at = NOW()
    `

    _, err := r.db.ExecContext(ctx, query, userID, objectKey, isUploaded)
    
    if err != nil {
        return fmt.Errorf("upsert failed: %w", err)
    }

    return nil
}