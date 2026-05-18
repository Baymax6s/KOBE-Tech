-- replies.kind を SMALLINT (0,1,2) から VARCHAR(20) ('comment','question','answer') に変更する

-- 1. 既存データを文字列に変換
UPDATE replies SET kind = CASE
    WHEN kind = 0 THEN 'comment'
    WHEN kind = 1 THEN 'question'
    WHEN kind = 2 THEN 'answer'
END::VARCHAR(20);

-- 2. カラム型を変更（CHECK制約も同時に更新）
ALTER TABLE replies
    ALTER COLUMN kind TYPE VARCHAR(20) USING kind::VARCHAR(20);

-- 3. 既存のCHECK制約を削除
ALTER TABLE replies
    DROP CONSTRAINT IF EXISTS chk_replies_kind;

-- 4. 新しいCHECK制約を追加
ALTER TABLE replies
    ADD CONSTRAINT chk_replies_kind
    CHECK (kind IN ('comment', 'question', 'answer'));
