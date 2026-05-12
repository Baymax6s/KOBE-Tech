package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/lib/pq"
)

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
