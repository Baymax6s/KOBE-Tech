package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
)

// PopularArticlesByTag はタグごとに「いいね数上位 limit 件」を返す。
//
// RANK() OVER (PARTITION BY tag ORDER BY いいね数 DESC) でタグ内順位を振り、
// CTE で 1 段挟んでから WHERE rank <= limit で上位だけに絞る
// （ウィンドウ関数は WHERE に直接書けないため）。
// 結果は tag_name, rank 順に並ぶので、タグ境界で TagRanking にまとめ直す。
func (r *Repository) PopularArticlesByTag(ctx context.Context, limit int) ([]article.TagRanking, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("article repository is not configured")
	}

	const query = `
		WITH ranked AS (
			SELECT
				t.id   AS tag_id,
				t.name AS tag_name,
				a.id   AS article_id,
				a.title,
				COUNT(l.id) AS like_count,
				RANK() OVER (
					PARTITION BY t.id
					ORDER BY COUNT(l.id) DESC, a.created_at DESC, a.id DESC
				) AS rank_in_tag
			FROM tags t
			JOIN article_tags at ON at.tag_id = t.id
			JOIN articles a      ON a.id = at.article_id
			LEFT JOIN likes l    ON l.article_id = a.id
			GROUP BY t.id, t.name, a.id, a.title
		)
		SELECT tag_id, tag_name, article_id, title, like_count, rank_in_tag
		FROM ranked
		WHERE rank_in_tag <= $1
		ORDER BY tag_name ASC, rank_in_tag ASC, article_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rankings := make([]article.TagRanking, 0)
	for rows.Next() {
		var (
			tagID   int64
			tagName string
			popular article.PopularArticle
		)
		if err := rows.Scan(
			&tagID,
			&tagName,
			&popular.ArticleID,
			&popular.Title,
			&popular.LikesCount,
			&popular.Rank,
		); err != nil {
			return nil, err
		}

		// クエリが tag 順に並んでいるので、直前と同じ tag なら追記、変われば新グループ。
		if n := len(rankings); n > 0 && rankings[n-1].TagID == tagID {
			rankings[n-1].Articles = append(rankings[n-1].Articles, popular)
			continue
		}
		rankings = append(rankings, article.TagRanking{
			TagID:    tagID,
			TagName:  tagName,
			Articles: []article.PopularArticle{popular},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rankings, nil
}
