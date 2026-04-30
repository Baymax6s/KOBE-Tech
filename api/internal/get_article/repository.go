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

func (r *Repository) FindByID(ctx context.Context, id int64) (Article, error) {
	if r == nil || r.db == nil {
		return Article{}, errors.New("article repository is not configured")
	}

	const query = `
		SELECT
			a.id,
			a.title,
			a.content,
			u.id,
			u.name,
			a.created_at,
			a.updated_at
		FROM articles a
		JOIN users u ON u.id = a.user_id
		WHERE a.id = $1
	`

	var article Article
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Author.ID,
		&article.Author.Name,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return Article{}, err
	}

	return article, nil
}
