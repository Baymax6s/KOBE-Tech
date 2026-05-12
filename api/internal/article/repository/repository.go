package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListArticles(ctx context.Context) ([]article.Article, error) {
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
			COALESCE(l.like_count, 0),
			COALESCE(tag_summary.tag_ids, ARRAY[]::integer[]),
			COALESCE(tag_summary.tag_names, ARRAY[]::text[])
		FROM articles a
		LEFT JOIN (
			SELECT article_id, COUNT(*) AS like_count FROM likes GROUP BY article_id
		) l ON l.article_id = a.id
		LEFT JOIN (
			SELECT
				article_tag.article_id,
				array_agg(t.id ORDER BY t.name, t.id) AS tag_ids,
				array_agg(t.name::text ORDER BY t.name, t.id) AS tag_names
			FROM article_tags article_tag
			JOIN tags t ON t.id = article_tag.tag_id
			GROUP BY article_tag.article_id
		) tag_summary ON tag_summary.article_id = a.id
		ORDER BY a.created_at DESC, a.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]article.Article, 0)
	for rows.Next() {
		var item article.Article
		var tagIDs []int64
		var tagNames []string
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.UserID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.LikesCount,
			pq.Array(&tagIDs),
			pq.Array(&tagNames),
		); err != nil {
			return nil, err
		}
		item.Tags = newTags(tagIDs, tagNames)

		articles = append(articles, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *Repository) FindArticleByID(ctx context.Context, id int64) (article.Article, error) {
	if r == nil || r.db == nil {
		return article.Article{}, errors.New("article repository is not configured")
	}

	const query = `
		SELECT
			a.id,
			a.title,
			a.content,
			u.id,
			u.name,
			a.created_at,
			a.updated_at,
			COALESCE((SELECT COUNT(*) FROM likes WHERE article_id = a.id), 0),
			COALESCE(tag_summary.tag_ids, ARRAY[]::integer[]),
			COALESCE(tag_summary.tag_names, ARRAY[]::text[])
		FROM articles a
		JOIN users u ON u.id = a.user_id
		LEFT JOIN LATERAL (
			SELECT
				array_agg(t.id ORDER BY t.name, t.id) AS tag_ids,
				array_agg(t.name::text ORDER BY t.name, t.id) AS tag_names
			FROM article_tags article_tag
			JOIN tags t ON t.id = article_tag.tag_id
			WHERE article_tag.article_id = a.id
		) tag_summary ON TRUE
		WHERE a.id = $1
	`

	var item article.Article
	var tagIDs []int64
	var tagNames []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.Author.ID,
		&item.Author.Name,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.LikesCount,
		pq.Array(&tagIDs),
		pq.Array(&tagNames),
	)
	if err != nil {
		return article.Article{}, err
	}
	item.Tags = newTags(tagIDs, tagNames)

	return item, nil
}

func (r *Repository) CreateArticle(ctx context.Context, title, content string, userID int64, tagNames []string) (article.Article, error) {
	if r == nil || r.db == nil {
		return article.Article{}, errors.New("post article repository is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return article.Article{}, err
	}
	defer tx.Rollback()

	const createArticleQuery = `
		INSERT INTO articles (title, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, title, content, user_id, created_at, updated_at
	`

	item := article.Article{
		Tags: make([]article.Tag, 0, len(tagNames)),
	}
	err = tx.QueryRowContext(ctx, createArticleQuery, title, content, userID).Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.UserID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return article.Article{}, err
	}

	for _, tagName := range tagNames {
		tag, err := upsertTag(ctx, tx, tagName)
		if err != nil {
			return article.Article{}, err
		}

		if err := attachTag(ctx, tx, item.ID, tag.ID); err != nil {
			return article.Article{}, err
		}

		item.Tags = append(item.Tags, tag)
	}

	if err := tx.Commit(); err != nil {
		return article.Article{}, err
	}

	return item, nil
}

func (r *Repository) ListTags(ctx context.Context) ([]article.Tag, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("tags repository is not configured")
	}

	const query = `
		SELECT id, name
		FROM tags
		ORDER BY name ASC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]article.Tag, 0)
	for rows.Next() {
		var tag article.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func newTags(ids []int64, names []string) []article.Tag {
	tags := make([]article.Tag, 0, len(ids))
	for i, id := range ids {
		if i >= len(names) {
			break
		}

		tags = append(tags, article.Tag{
			ID:   id,
			Name: names[i],
		})
	}

	return tags
}

func upsertTag(ctx context.Context, tx *sql.Tx, name string) (article.Tag, error) {
	const selectQuery = `
		SELECT id, name
		FROM tags
		WHERE name = $1
	`

	var tag article.Tag
	err := tx.QueryRowContext(ctx, selectQuery, name).Scan(&tag.ID, &tag.Name)
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return article.Tag{}, err
	}

	const insertQuery = `
		INSERT INTO tags (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		ON CONFLICT (name) DO NOTHING
		RETURNING id, name
	`

	err = tx.QueryRowContext(ctx, insertQuery, name).Scan(&tag.ID, &tag.Name)
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return article.Tag{}, err
	}

	err = tx.QueryRowContext(ctx, selectQuery, name).Scan(&tag.ID, &tag.Name)
	if err != nil {
		return article.Tag{}, err
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
