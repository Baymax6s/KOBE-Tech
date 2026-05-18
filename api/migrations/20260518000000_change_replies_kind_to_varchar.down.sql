-- replies.kind を VARCHAR(20) から SMALLINT に戻す

-- 1. 既存データを数値に変換
UPDATE replies SET kind = CASE
    WHEN kind = 'comment' THEN '0'
    WHEN kind = 'question' THEN '1'
    WHEN kind = 'answer' THEN '2'
END::VARCHAR(20);

-- 2. カラム型を変更（CHECK制約も同時に更新）
ALTER TABLE replies
    ALTER COLUMN kind TYPE SMALLINT USING kind::SMALLINT;

-- 3. 既存のCHECK制約を削除
ALTER TABLE replies
    DROP CONSTRAINT IF EXISTS chk_replies_kind;

-- 4. 元のCHECK制約を追加
ALTER TABLE replies
    ADD CONSTRAINT chk_replies_kind
    CHECK (kind IN (0, 1, 2));
