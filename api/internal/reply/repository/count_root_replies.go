package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
)

// CountRootByKind は記事のルート返信（parent_id IS NULL）を kind 別に集計する。
// ルート返信は comment / question のみ（answer は必ず親を持つ）なので、All = comment + question。
func (r *Repository) CountRootByKind(ctx context.Context, articleID int64) (reply.RootKindCounts, error) {
	if r == nil || r.db == nil {
		return reply.RootKindCounts{}, errors.New("reply repository is not configured")
	}

	const query = `
		SELECT kind, COUNT(*)
		FROM replies
		WHERE article_id = $1 AND parent_id IS NULL
		GROUP BY kind
	`

	rows, err := r.db.QueryContext(ctx, query, articleID)
	if err != nil {
		return reply.RootKindCounts{}, err
	}
	defer rows.Close()

	var counts reply.RootKindCounts
	for rows.Next() {
		var kind reply.Kind
		var n int64
		if err := rows.Scan(&kind, &n); err != nil {
			return reply.RootKindCounts{}, err
		}
		counts.All += n
		switch kind {
		case reply.KindComment:
			counts.Comment = n
		case reply.KindQuestion:
			counts.Question = n
		}
	}
	if err := rows.Err(); err != nil {
		return reply.RootKindCounts{}, err
	}

	return counts, nil
}
