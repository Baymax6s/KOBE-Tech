package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
)

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
