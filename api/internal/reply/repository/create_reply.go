package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
	"github.com/lib/pq"
)

func (r *Repository) Create(ctx context.Context, articleID int64, userID int64, parentID *int64, kind reply.Kind, body string) (reply.Reply, error) {
	if r == nil || r.db == nil {
		return reply.Reply{}, errors.New("reply repository is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return reply.Reply{}, err
	}
	defer tx.Rollback()

	if err := ensureArticleExists(ctx, tx, articleID); err != nil {
		return reply.Reply{}, err
	}

	if parentID == nil {
		// ルート投稿は comment か question のみ許可（answer は必ず親が必要）。
		if kind != reply.KindComment && kind != reply.KindQuestion {
			return reply.Reply{}, ErrInvalidRootKind
		}
	} else {
		parentKind, parentArticleID, err := fetchParent(ctx, tx, *parentID)
		if err != nil {
			return reply.Reply{}, err
		}
		if parentArticleID != articleID {
			return reply.Reply{}, ErrParentMismatch
		}
		// 親 kind ごとに子 kind を一意に決める:
		//   comment への返信は comment、question/answer への返信は answer。
		// これによりスレッドの種別がフォーム送信のタイミングで競合しないことを保証する。
		var expected reply.Kind
		if parentKind == reply.KindComment {
			expected = reply.KindComment
		} else {
			expected = reply.KindAnswer
		}
		if kind != expected {
			return reply.Reply{}, ErrInvalidParent
		}
	}

	const insertQuery = `
		INSERT INTO replies (article_id, user_id, parent_id, kind, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, article_id, parent_id, kind, content, user_id, created_at, updated_at
	`

	var item reply.Reply
	var parentVal sql.NullInt64
	err = tx.QueryRowContext(ctx, insertQuery, articleID, userID, nullableInt64(parentID), string(kind), body).Scan(
		&item.ID,
		&item.ArticleID,
		&parentVal,
		&item.Kind,
		&item.Body,
		&item.UserID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			switch pqErr.Constraint {
			case "fk_replies_article":
				return reply.Reply{}, ErrArticleNotFound
			case "fk_replies_parent":
				return reply.Reply{}, ErrParentNotFound
			}
		}
		return reply.Reply{}, err
	}
	if parentVal.Valid {
		v := parentVal.Int64
		item.ParentID = &v
	}

	const userQuery = `SELECT name FROM users WHERE id = $1`
	if err := tx.QueryRowContext(ctx, userQuery, userID).Scan(&item.UserName); err != nil {
		return reply.Reply{}, err
	}

	if err := tx.Commit(); err != nil {
		return reply.Reply{}, err
	}

	return item, nil
}

func ensureArticleExists(ctx context.Context, tx *sql.Tx, articleID int64) error {
	const query = `SELECT 1 FROM articles WHERE id = $1`
	var exists int
	err := tx.QueryRowContext(ctx, query, articleID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrArticleNotFound
	}
	return err
}

func fetchParent(ctx context.Context, tx *sql.Tx, parentID int64) (reply.Kind, int64, error) {
	const query = `SELECT kind, article_id FROM replies WHERE id = $1`
	var kind reply.Kind
	var articleID int64
	err := tx.QueryRowContext(ctx, query, parentID).Scan(&kind, &articleID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, ErrParentNotFound
	}
	if err != nil {
		return "", 0, err
	}
	return kind, articleID, nil
}

func nullableInt64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}
