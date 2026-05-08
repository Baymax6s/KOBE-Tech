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

func (r *Repository) Create(ctx context.Context, title, content string, userID int64, tagNames []string) (Article, error) {
	if r == nil || r.db == nil {
		return Article{}, errors.New("post article repository is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Article{}, err
	}
	defer tx.Rollback()

	const createArticleQuery = `
		INSERT INTO articles (title, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, title, content, user_id, created_at, updated_at
	`

	article := Article{
		Tags: make([]Tag, 0, len(tagNames)),
	}
	err = tx.QueryRowContext(ctx, createArticleQuery, title, content, userID).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return Article{}, err
	}

	for _, tagName := range tagNames {
		tag, err := upsertTag(ctx, tx, tagName)
		if err != nil {
			return Article{}, err
		}

		if err := attachTag(ctx, tx, article.ID, tag.ID); err != nil {
			return Article{}, err
		}

		article.Tags = append(article.Tags, tag)
	}

	if err := tx.Commit(); err != nil {
		return Article{}, err
	}

	return article, nil
}

func upsertTag(ctx context.Context, tx *sql.Tx, name string) (Tag, error) {
	const insertQuery = `
		INSERT INTO tags (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		ON CONFLICT (name) DO NOTHING
		RETURNING id, name
	`

	var tag Tag
	err := tx.QueryRowContext(ctx, insertQuery, name).Scan(&tag.ID, &tag.Name)
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Tag{}, err
	}

	const selectQuery = `
		SELECT id, name
		FROM tags
		WHERE name = $1
	`
	err = tx.QueryRowContext(ctx, selectQuery, name).Scan(&tag.ID, &tag.Name)
	if err != nil {
		return Tag{}, err
	}

	return tag, nil
}

func attachTag(ctx context.Context, tx *sql.Tx, articleID, tagID int64) error {
	const query = `
		INSERT INTO article_tags (article_id, tag_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
	`

	_, err := tx.ExecContext(ctx, query, articleID, tagID)
	return err
}
