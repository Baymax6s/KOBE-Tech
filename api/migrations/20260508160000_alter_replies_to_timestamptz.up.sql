-- replies.created_at / updated_at を他テーブルと揃えて TIMESTAMPTZ に変更する。
-- 既存値はセッション TZ 由来のため、同じ TZ で再解釈してから TIMESTAMPTZ にする。
ALTER TABLE replies
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
        USING created_at AT TIME ZONE current_setting('TimeZone'),
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ
        USING updated_at AT TIME ZONE current_setting('TimeZone');
