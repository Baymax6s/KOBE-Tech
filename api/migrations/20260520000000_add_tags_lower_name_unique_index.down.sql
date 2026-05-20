-- 削除した tags_name_key を復活させてから式インデックスを落とす。
ALTER TABLE tags ADD CONSTRAINT tags_name_key UNIQUE (name);
DROP INDEX IF EXISTS uniq_tags_lower_name;
