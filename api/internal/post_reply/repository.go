package reply

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	errArticleNotFound = errors.New("article not found")
	errParentNotFound  = errors.New("parent reply not found")
	errParentMismatch  = errors.New("parent reply does not belong to the article")
	errInvalidParent   = errors.New("parent reply kind is not allowed for this reply kind")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, articleID int64, userID int64, parentID *int64, kind Kind, body string) (Reply, error) {
	if r == nil || r.db == nil {
		return Reply{}, errors.New("post reply repository is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Reply{}, err
	}
	defer tx.Rollback()

	if err := ensureArticleExists(ctx, tx, articleID); err != nil {
		return Reply{}, err
	}

	if parentID != nil {
		parentKind, parentArticleID, err := fetchParent(ctx, tx, *parentID)
		if err != nil {
			return Reply{}, err
		}
		if parentArticleID != articleID {
			return Reply{}, errParentMismatch
		}
		// 今回はコメントのみ。親が comment のときに限り comment を返信できる。
		if parentKind != KindComment || kind != KindComment {
			return Reply{}, errInvalidParent
		}
	}

	const insertQuery = `
		INSERT INTO replies (article_id, user_id, parent_id, kind, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, article_id, parent_id, kind, content, user_id, created_at, updated_at
	`

	var reply Reply
	var kindVal int16
	var parentVal sql.NullInt64
	err = tx.QueryRowContext(ctx, insertQuery, articleID, userID, nullableInt64(parentID), int16(kind), body).Scan(
		&reply.ID,
		&reply.ArticleID,
		&parentVal,
		&kindVal,
		&reply.Body,
		&reply.UserID,
		&reply.CreatedAt,
		&reply.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			switch pqErr.Constraint {
			case "fk_replies_article":
				return Reply{}, errArticleNotFound
			case "fk_replies_parent":
				return Reply{}, errParentNotFound
			}
		}
		return Reply{}, err
	}
	if parentVal.Valid {
		v := parentVal.Int64
		reply.ParentID = &v
	}
	reply.Kind = Kind(kindVal)

	const userQuery = `SELECT name FROM users WHERE id = $1`
	if err := tx.QueryRowContext(ctx, userQuery, userID).Scan(&reply.UserName); err != nil {
		return Reply{}, err
	}

	if err := tx.Commit(); err != nil {
		return Reply{}, err
	}

	return reply, nil
}

func ensureArticleExists(ctx context.Context, tx *sql.Tx, articleID int64) error {
	const query = `SELECT 1 FROM articles WHERE id = $1`
	var exists int
	err := tx.QueryRowContext(ctx, query, articleID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return errArticleNotFound
	}
	return err
}

func fetchParent(ctx context.Context, tx *sql.Tx, parentID int64) (Kind, int64, error) {
	const query = `SELECT kind, article_id FROM replies WHERE id = $1`
	var kindVal int16
	var articleID int64
	err := tx.QueryRowContext(ctx, query, parentID).Scan(&kindVal, &articleID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, errParentNotFound
	}
	if err != nil {
		return 0, 0, err
	}
	return Kind(kindVal), articleID, nil
}

func nullableInt64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}
