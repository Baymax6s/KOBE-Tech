-- 1. users に bio を戻す
ALTER TABLE users ADD COLUMN bio TEXT;

-- 2. user_profiles のデータを users に戻す
UPDATE users
SET bio = user_profiles.bio
FROM user_profiles
WHERE user_profiles.user_id = users.id;

-- 3. テーブル削除
DROP TABLE IF EXISTS user_profiles;