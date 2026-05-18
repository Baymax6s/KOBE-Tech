-- replies.kind を SMALLINT (0,1,2) から VARCHAR(20) ('comment','question','answer') に変更する

-- 1. 既存のCHECK制約を削除（型変更前に必要）
ALTER TABLE replies
    DROP CONSTRAINT IF EXISTS chk_replies_kind;

-- 2. カラム型の変更と既存データの値変換を同時に実施
ALTER TABLE replies
    ALTER COLUMN kind TYPE VARCHAR(20) USING (
        CASE
            WHEN kind = 0 THEN 'comment'
            WHEN kind = 1 THEN 'question'
            WHEN kind = 2 THEN 'answer'
        END
    );

-- 3. 新しいCHECK制約を追加
ALTER TABLE replies
    ADD CONSTRAINT chk_replies_kind
    CHECK (kind IN ('comment', 'question', 'answer'));
