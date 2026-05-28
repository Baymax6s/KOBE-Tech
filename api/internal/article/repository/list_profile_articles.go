package repository

import (
	"context"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/article"
)

// ListArticlesByAuthor は authorID が投稿した記事を新しい順に返す。
// viewerID は閲覧者で、各記事の liked_by_me 判定に使う（authorID とは別概念）。
func (r *Repository) ListArticlesByAuthor(ctx context.Context, authorID int64, viewerID int64) ([]article.Article, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("article repository is not configured")
	}

	query := articleListSelectFrom + `
		WHERE a.user_id = $2
		ORDER BY a.created_at DESC, a.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, viewerID, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleList(rows)
}

// ListBestAnswerArticles は userID の回答がベストアンサーに選ばれた記事を新しい順に返す。
// viewerID は閲覧者で、各記事の liked_by_me 判定に使う（userID とは別概念）。
func (r *Repository) ListBestAnswerArticles(ctx context.Context, userID int64, viewerID int64) ([]article.Article, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("article repository is not configured")
	}

	query := articleListSelectFrom + `
		WHERE EXISTS (
			SELECT 1 FROM replies r
			WHERE r.article_id = a.id AND r.user_id = $2 AND r.is_best = TRUE
		)
		ORDER BY a.created_at DESC, a.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, viewerID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleList(rows)
}

// ListLikedArticles は likerID がいいねした記事を「いいねした順」で返す。
// viewerID は閲覧者で、各記事の liked_by_me 判定に使う（likerID とは別概念）。
func (r *Repository) ListLikedArticles(ctx context.Context, likerID int64, viewerID int64) ([]article.Article, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("article repository is not configured")
	}

	query := articleListSelectFrom + `
		JOIN likes liker_like
			ON liker_like.article_id = a.id AND liker_like.user_id = $2
		ORDER BY liker_like.created_at DESC, a.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, viewerID, likerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleList(rows)
}
