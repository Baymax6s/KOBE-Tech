package post_reply

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// =======================
// FindByID（追加）
// =======================
func (r *Repository) FindByID(ctx context.Context, id int64) (Reply, error) {
	if r == nil || r.db == nil {
		return Reply{}, errors.New("post reply repository not configured")
	}

	const query = `
	SELECT
		id,
		article_id,
		user_id,
		content,
		parent_id,
		kind,
		is_best,
		created_at,
		updated_at
	FROM replies
	WHERE id = $1
	`

	var rp Reply

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rp.ID,
		&rp.ArticleID,
		&rp.UserID,
		&rp.Content,
		&rp.ParentID,
		&rp.Kind,
		&rp.IsBest,
		&rp.CreatedAt,
		&rp.UpdatedAt,
	)

	return rp, err
}

// =======================
// Create（既存）
// =======================
func (r *Repository) Create(
	ctx context.Context,
	articleID, userID int64,
	content string,
	parentID *int64,
	kind int,
) (Reply, error) {

	if r == nil || r.db == nil {
		return Reply{}, errors.New("post reply repository not configured")
	}

	const query = `
	WITH inserted AS (
		INSERT INTO replies (
			article_id,
			user_id,
			content,
			parent_id,
			kind
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			article_id,
			user_id,
			content,
			parent_id,
			kind,
			is_best,
			created_at,
			updated_at
	)
	SELECT
		i.id,
		i.article_id,
		i.user_id,
		i.content,
		i.parent_id,
		i.kind,
		i.is_best,
		u.name,
		i.created_at,
		i.updated_at
	FROM inserted i
	JOIN users u ON u.id = i.user_id
	`

	var rp Reply

	err := r.db.QueryRowContext(
		ctx,
		query,
		articleID,
		userID,
		content,
		parentID,
		kind,
	).Scan(
		&rp.ID,
		&rp.ArticleID,
		&rp.UserID,
		&rp.Content,
		&rp.ParentID,
		&rp.Kind,
		&rp.IsBest,
		&rp.UserName,
		&rp.CreatedAt,
		&rp.UpdatedAt,
	)

	return rp, err
}