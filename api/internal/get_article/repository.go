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
			COALESCE(
				array_agg(t.id ORDER BY t.name, t.id) FILTER (WHERE t.id IS NOT NULL),
				ARRAY[]::integer[]
			),
			COALESCE(
				array_agg(t.name::text ORDER BY t.name, t.id) FILTER (WHERE t.id IS NOT NULL),
				ARRAY[]::text[]
			)
		FROM articles a
		JOIN users u ON u.id = a.user_id
		LEFT JOIN article_tags article_tag ON article_tag.article_id = a.id
		LEFT JOIN tags t ON t.id = article_tag.tag_id
		WHERE a.id = $1
		GROUP BY
			a.id,
			a.title,
			a.content,
			u.id,
			u.name,
			a.created_at,
			a.updated_at
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
