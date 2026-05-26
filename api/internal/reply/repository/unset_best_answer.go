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

	// 1. 指定された reply が存在し、answer かつベストアンサーであることを確認する。
	const fetchReplyQuery = `
		SELECT kind, is_best
		FROM replies
		WHERE id = $1
	`
	var kind reply.Kind
	var isBest bool
	err = tx.QueryRowContext(ctx, fetchReplyQuery, replyID).Scan(&kind, &isBest)
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

	// 2. ルートの質問者を特定して権限チェックを行う。
	const findRootQuery = `
		WITH RECURSIVE root_path AS (
			SELECT id, parent_id, kind, user_id
			FROM replies
			WHERE id = $1
			UNION ALL
			SELECT r.id, r.parent_id, r.kind, r.user_id
			FROM replies r
			JOIN root_path rp ON r.id = rp.parent_id
		)
		SELECT kind, user_id FROM root_path WHERE parent_id IS NULL
	`
	var rootKind reply.Kind
	var questionUserID int64
	err = tx.QueryRowContext(ctx, findRootQuery, replyID).Scan(&rootKind, &questionUserID)
	if err != nil {
		return err
	}

	if rootKind != reply.KindQuestion {
		return ErrNotAnswer
	}

	if questionUserID != userID {
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
