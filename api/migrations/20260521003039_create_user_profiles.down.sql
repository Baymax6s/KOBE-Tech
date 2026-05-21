-- 1. 削除してしまった users テーブルの bio カラムを復活させる
ALTER TABLE users ADD COLUMN bio TEXT;

-- 2. 作成した新テーブルを削除する
DROP TABLE IF EXISTS user_profiles;