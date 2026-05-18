package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
)

func (r *Repository) ListByArticleID(ctx context.Context, articleID int64) ([]reply.Reply, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("reply repository is not configured")
	}

	const query = `
		SELECT
			r.id,
			r.article_id,
			r.parent_id,
			r.kind,
			r.content,
			r.user_id,
			u.name,
			r.created_at,
			r.updated_at
		FROM replies r
		JOIN users u ON u.id = r.user_id
		WHERE r.article_id = $1
		ORDER BY r.created_at ASC, r.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	replies := make([]reply.Reply, 0)
	for rows.Next() {
		var item reply.Reply
		var parent sql.NullInt64
		if err := rows.Scan(
			&item.ID,
			&item.ArticleID,
			&parent,
			&item.Kind,
			&item.Body,
			&item.UserID,
			&item.UserName,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if parent.Valid {
			v := parent.Int64
			item.ParentID = &v
		}
		replies = append(replies, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return replies, nil
}
