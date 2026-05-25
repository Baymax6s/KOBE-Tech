package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
)

func (r *Repository) UnsetBestAnswer(ctx context.Context, replyID, userID int64) error {
	if r == nil || r.db == nil {
		return errors.New("reply repository is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 回答(a)と、その親である質問(q)を replies の自己結合で一度に引く。
	// parent_id が NULL（ルート投稿）でも回答自身の行は返したいので LEFT JOIN。
	const fetchReplyQuery = `
		SELECT a.kind, a.parent_id, a.is_best, q.kind, q.user_id
		FROM replies a
		LEFT JOIN replies q ON q.id = a.parent_id
		WHERE a.id = $1
	`
	var kind reply.Kind
	var parentID sql.NullInt64
	var isBest bool
	var parentKind sql.NullString
	var questionUserID sql.NullInt64
	err = tx.QueryRowContext(ctx, fetchReplyQuery, replyID).Scan(&kind, &parentID, &isBest, &parentKind, &questionUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrReplyNotFound
	}
	if err != nil {
		return err
	}
	if kind != reply.KindAnswer {
		return ErrNotAnswer
	}
	if !isBest {
		return ErrNotBestAnswer
	}
	if !parentID.Valid {
		return ErrNotAnswer
	}
	if !parentKind.Valid {
		return ErrParentNotFound
	}
	if reply.Kind(parentKind.String) != reply.KindQuestion {
		return ErrNotAnswer
	}
	if questionUserID.Int64 != userID {
		return ErrNotQuestionAuthor
	}

	const updateQuery = `UPDATE replies SET is_best = FALSE, updated_at = NOW() WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateQuery, replyID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
