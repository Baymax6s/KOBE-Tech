-- 4人分のプロフィールとアバター画像のシードデータ
INSERT INTO user_profiles (user_id, object_key, bio, is_uploaded) VALUES
(1, 'avatar/1.jpg', '今はgitの勉強を頑張っています！', TRUE), 
(2, 'avatar/2.jpg', 'こんにちは！主にテック系の記事を投稿します。', TRUE),
(3, 'avatar/3.jpg', 'Go言語とクリーンアーキテクチャに興味があります。', TRUE),
(4, 'avatar/4.png', '新米エンジニアです。よろしくお願いします！', TRUE)
ON CONFLICT (user_id) DO UPDATE 
SET object_key = EXCLUDED.object_key, 
    bio = EXCLUDED.bio, 
    is_uploaded = EXCLUDED.is_uploaded,
    updated_at = CURRENT_TIMESTAMP;