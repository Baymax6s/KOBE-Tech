-- replies.kind を VARCHAR(20) から SMALLINT に戻す

-- 1. VARCHAR用のCHECK制約を先に削除
ALTER TABLE replies
    DROP CONSTRAINT IF EXISTS chk_replies_kind;

-- 2. カラム型の変更と既存データの値変換を同時に実施
ALTER TABLE replies
    ALTER COLUMN kind TYPE SMALLINT USING (
        CASE
            WHEN kind = 'comment' THEN 0
            WHEN kind = 'question' THEN 1
            WHEN kind = 'answer' THEN 2
        END
    );

-- 3. 元のCHECK制約を追加
ALTER TABLE replies
    ADD CONSTRAINT chk_replies_kind
    CHECK (kind IN (0, 1, 2));
