package list_replies

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

// 一覧取得（記事ID単位）
func (r *Repository) FindByArticleID(ctx context.Context, articleID int64) ([]Reply, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("list replies repository not configured")
	}

	const query = `
	SELECT
		r.id,
		r.article_id,
		r.user_id,
		r.content,
		r.parent_id,
		r.kind,
		r.is_best,
		u.name,
		r.created_at,
		r.updated_at
	FROM replies r
	JOIN users u ON u.id = r.user_id
	WHERE r.article_id = $1
	ORDER BY r.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Reply

	for rows.Next() {
		var rp Reply

		if err := rows.Scan(
			&rp.ID,
			&rp.ArticleID,
			&rp.UserID,
			&rp.Content,
			&rp.ParentID,
			&rp.Kind,
			&rp.IsBest,
			&rp.UserName,
			&rp.CreatedAt,
			&rp.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, rp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}