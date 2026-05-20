package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
	"github.com/lib/pq"
)

func (r *Repository) ListArticles(ctx context.Context, userID int64, tagNames []string) ([]article.Article, error) {
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
			COALESCE(tag_summary.tag_names, ARRAY[]::text[]),
			EXISTS(SELECT 1 FROM likes WHERE article_id = a.id AND user_id = $1)
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
	`

	args := []any{userID}
	if len(tagNames) > 0 {
		// 指定された tag をすべて持つ記事のみに絞る（AND セマンティクス）。
		// 大文字小文字は区別しないので比較は LOWER(name) で行い、tagNames も lowercase 済みを渡す。
		// HAVING COUNT(DISTINCT LOWER(t.name)) = $3 で「一致した tag の種類数 == 要求された tag 数」を担保。
		query += `
		WHERE a.id IN (
			SELECT article_tag.article_id
			FROM article_tags article_tag
			JOIN tags t ON t.id = article_tag.tag_id
			WHERE LOWER(t.name) = ANY($2::text[])
			GROUP BY article_tag.article_id
			HAVING COUNT(DISTINCT LOWER(t.name)) = $3
		)
		`
		args = append(args, pq.Array(tagNames), len(tagNames))
	}
	query += ` ORDER BY a.created_at DESC, a.id DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
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
			&item.LikedByMe,
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
