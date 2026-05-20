package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
)

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

func upsertTag(ctx context.Context, tx *sql.Tx, name string) (article.Tag, error) {
	// "Vue" と "vue" を同一タグとして扱うため、検索も衝突判定も LOWER(name) で行う。
	// 既存行が見つかった場合は最初に登録された原文ケース（tag.Name）をそのまま返す。
	const selectQuery = `
		SELECT id, name
		FROM tags
		WHERE LOWER(name) = LOWER($1)
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
		ON CONFLICT (LOWER(name)) DO NOTHING
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
