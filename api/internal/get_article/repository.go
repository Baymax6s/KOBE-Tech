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

	var article Article
	var tagIDs []int64
	var tagNames []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Author.ID,
		&article.Author.Name,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.LikesCount,
		pq.Array(&tagIDs),
		pq.Array(&tagNames),
	)
	if err != nil {
		return Article{}, err
	}
	article.Tags = newTags(tagIDs, tagNames)

	return article, nil
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
