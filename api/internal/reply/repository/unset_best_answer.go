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

	// 1. 指定された reply が存在し、ルートの質問者を特定して権限チェックを行う。
	// クエリ数を削減するため、Recursive CTE でターゲットの情報を引き継ぎつつルート（parent_id IS NULL）を特定する。
	const findRootQuery = `
		WITH RECURSIVE root_path AS (
			SELECT id, parent_id, kind, user_id, kind as target_kind, is_best as target_is_best
			FROM replies
			WHERE id = $1
			UNION ALL
			SELECT r.id, r.parent_id, r.kind, r.user_id, rp.target_kind, rp.target_is_best
			FROM replies r
			JOIN root_path rp ON r.id = rp.parent_id
		)
		SELECT target_kind, target_is_best, kind, user_id FROM root_path WHERE parent_id IS NULL
	`

	var targetKind reply.Kind
	var targetIsBest bool
	var rootKind reply.Kind
	var questionUserID int64
	err = tx.QueryRowContext(ctx, findRootQuery, replyID).Scan(&targetKind, &targetIsBest, &rootKind, &questionUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrReplyNotFound
	}
	if err != nil {
		return err
	}

	// ターゲットの状態チェック
	if targetKind != reply.KindAnswer {
		return ErrNotAnswer
	}
	if !targetIsBest {
		return ErrNotBestAnswer
	}

	// ルートの整合性と権限チェック
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
