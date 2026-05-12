package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
	"github.com/lib/pq"
)

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrParentNotFound  = errors.New("parent reply not found")
	ErrParentMismatch  = errors.New("parent reply does not belong to the article")
	ErrInvalidParent   = errors.New("parent reply kind is not allowed for this reply kind")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByArticleID(ctx context.Context, articleID int64) ([]reply.Reply, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("reply repository is not configured")
	}

	const query = `
		SELECT
			r.id,
			r.article_id,
			r.parent_id,
			r.kind,
			r.content,
			r.user_id,
			u.name,
			r.created_at,
			r.updated_at
		FROM replies r
		JOIN users u ON u.id = r.user_id
		WHERE r.article_id = $1
		ORDER BY r.created_at ASC, r.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	replies := make([]reply.Reply, 0)
	for rows.Next() {
		var item reply.Reply
		var parent sql.NullInt64
		var kindVal int16
		if err := rows.Scan(
			&item.ID,
			&item.ArticleID,
			&parent,
			&kindVal,
			&item.Body,
			&item.UserID,
			&item.UserName,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if parent.Valid {
			v := parent.Int64
			item.ParentID = &v
		}
		item.Kind = reply.Kind(kindVal)
		replies = append(replies, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return replies, nil
}

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

	if parentID != nil {
		parentKind, parentArticleID, err := fetchParent(ctx, tx, *parentID)
		if err != nil {
			return reply.Reply{}, err
		}
		if parentArticleID != articleID {
			return reply.Reply{}, ErrParentMismatch
		}
		// 今回はコメントのみ。親が comment のときに限り comment を返信できる。
		if parentKind != reply.KindComment || kind != reply.KindComment {
			return reply.Reply{}, ErrInvalidParent
		}
	}

	const insertQuery = `
		INSERT INTO replies (article_id, user_id, parent_id, kind, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, article_id, parent_id, kind, content, user_id, created_at, updated_at
	`

	var item reply.Reply
	var kindVal int16
	var parentVal sql.NullInt64
	err = tx.QueryRowContext(ctx, insertQuery, articleID, userID, nullableInt64(parentID), int16(kind), body).Scan(
		&item.ID,
		&item.ArticleID,
		&parentVal,
		&kindVal,
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
	item.Kind = reply.Kind(kindVal)

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
	var kindVal int16
	var articleID int64
	err := tx.QueryRowContext(ctx, query, parentID).Scan(&kindVal, &articleID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, ErrParentNotFound
	}
	if err != nil {
		return 0, 0, err
	}
	return reply.Kind(kindVal), articleID, nil
}

func nullableInt64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}
