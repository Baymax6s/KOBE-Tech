package article

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
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
			COALESCE(l.like_count, 0),
			COALESCE(
				array_agg(t.id ORDER BY t.name, t.id) FILTER (WHERE t.id IS NOT NULL),
				ARRAY[]::integer[]
			),
			COALESCE(
				array_agg(t.name::text ORDER BY t.name, t.id) FILTER (WHERE t.id IS NOT NULL),
				ARRAY[]::text[]
			)
		FROM articles a
		LEFT JOIN (
			SELECT article_id, COUNT(*) AS like_count FROM likes GROUP BY article_id
		) l ON l.article_id = a.id
		LEFT JOIN article_tags article_tag ON article_tag.article_id = a.id
		LEFT JOIN tags t ON t.id = article_tag.tag_id
		GROUP BY
			a.id,
			a.title,
			a.content,
			a.user_id,
			a.created_at,
			a.updated_at,
			l.like_count
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
		var tagIDs []int64
		var tagNames []string
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.LikesCount,
			pq.Array(&tagIDs),
			pq.Array(&tagNames),
		); err != nil {
			return nil, err
		}
		article.Tags = newTags(tagIDs, tagNames)

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func newTags(ids []int64, names []string) []Tag {
	tags := make([]Tag, 0, len(ids))
	for i, id := range ids {
		if i >= len(names) {
			break
		}

		tags = append(tags, Tag{
			ID:   id,
			Name: names[i],
		})
	}

	return tags
}
