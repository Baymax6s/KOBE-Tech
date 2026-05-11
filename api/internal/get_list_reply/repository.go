package reply

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

func (r *Repository) ListByArticleID(ctx context.Context, articleID int64) ([]Reply, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("get list reply repository is not configured")
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

	replies := make([]Reply, 0)
	for rows.Next() {
		var reply Reply
		var parent sql.NullInt64
		var kindVal int16
		if err := rows.Scan(
			&reply.ID,
			&reply.ArticleID,
			&parent,
			&kindVal,
			&reply.Body,
			&reply.UserID,
			&reply.UserName,
			&reply.CreatedAt,
			&reply.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if parent.Valid {
			v := parent.Int64
			reply.ParentID = &v
		}
		reply.Kind = Kind(kindVal)
		replies = append(replies, reply)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return replies, nil
}
