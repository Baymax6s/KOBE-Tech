package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Baymax6s/KOBE-Tech/api/internal/reply"
)

func (r *Repository) SetBestAnswer(ctx context.Context, replyID, userID int64) error {
	if r == nil || r.db == nil {
		return errors.New("reply repository is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 指定された reply が存在し、answer であることを確認し、同時にルートの質問者 ID を取得する。
	// ネスト制限が最大 2 なので、再帰的に遡るか、直接 root を特定する。
	// ここでは汎用性のために Recursive CTE を使用してルート投稿（parent_id IS NULL）を見つける。
	const findRootQuery = `
		WITH RECURSIVE root_path AS (
			SELECT id, parent_id, kind, user_id, 1 as depth
			FROM replies
			WHERE id = $1
			UNION ALL
			SELECT r.id, r.parent_id, r.kind, r.user_id, rp.depth + 1
			FROM replies r
			JOIN root_path rp ON r.id = rp.parent_id
		)
		SELECT id, kind, user_id FROM root_path WHERE parent_id IS NULL
	`

	var rootID int64
	var rootKind reply.Kind
	var questionUserID int64

	// まず指定された reply 自体の kind をチェック
	const fetchKindQuery = `SELECT kind FROM replies WHERE id = $1`
	var targetKind reply.Kind
	err = tx.QueryRowContext(ctx, fetchKindQuery, replyID).Scan(&targetKind)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrReplyNotFound
	}
	if err != nil {
		return err
	}
	if targetKind != reply.KindAnswer {
		return ErrNotAnswer
	}

	// ルートを特定して権限チェック
	err = tx.QueryRowContext(ctx, findRootQuery, replyID).Scan(&rootID, &rootKind, &questionUserID)
	if errors.Is(err, sql.ErrNoRows) {
		// 自分がルートの場合は親がいないのでここに来る可能性があるが、KindAnswer なのであり得ないはず。
		return ErrNotAnswer
	}
	if err != nil {
		return err
	}

	if rootKind != reply.KindQuestion {
		return ErrNotAnswer // 質問スレッド以外でのベストアンサーは不可
	}

	if questionUserID != userID {
		return ErrNotQuestionAuthor
	}

	// 2. スレッド内のいずれかの返信にすでにベストアンサーが設定されていないか確認する。
	// ルート質問 ID を起点に、その配下の全子孫（再帰）の中で is_best = TRUE のものを探す。
	const checkThreadBestQuery = `
		WITH RECURSIVE thread_members AS (
			SELECT id FROM replies WHERE id = $1
			UNION ALL
			SELECT r.id FROM replies r
			JOIN thread_members tm ON r.parent_id = tm.id
		)
		SELECT EXISTS (
			SELECT 1 FROM replies WHERE id IN (SELECT id FROM thread_members) AND is_best = TRUE
		)
	`
	var hasBest bool
	err = tx.QueryRowContext(ctx, checkThreadBestQuery, rootID).Scan(&hasBest)
	if err != nil {
		return err
	}
	if hasBest {
		return ErrBestAnswerAlreadySet
	}

	const updateQuery = `UPDATE replies SET is_best = TRUE, updated_at = NOW() WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateQuery, replyID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
