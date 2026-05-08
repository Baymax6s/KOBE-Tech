package article

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

func (r *Repository) List(ctx context.Context) ([]Article, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("article repository is not configured")
	}
	query := `
		SELECT
			a.id,
			a.title,
			a.content,
			a.user_id,
			a.created_at,
			a.updated_at,
			COALESCE(l.like_count, 0)
		FROM articles a
		LEFT JOIN (
			SELECT article_id, COUNT(*) AS like_count FROM likes GROUP BY article_id
		) l ON l.article_id = a.id
		ORDER BY a.created_at DESC, a.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.LikesCount,
		); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
