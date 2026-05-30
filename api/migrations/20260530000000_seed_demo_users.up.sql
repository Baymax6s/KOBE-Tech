-- デモ用に「実在しそうな名前」のユーザーを追加し、プロフィールの各状態を網羅する。
-- 既存の admin / 田中太郎 / 山田花子 / 佐藤次郎 (id 1-4) は触らない。
--
-- 追加するプロフィールの状態（デモで見せたいバリエーション）:
--   - アバター画像あり (is_uploaded = TRUE)        … 鈴木健一 / 高橋美咲 / 渡辺翔太
--   - アバターなし・自己紹介あり (is_uploaded=FALSE) … 伊藤さくら / 中村大輔
--   - プロフィール行そのものが無い最小状態          … 小林優子（投稿せず「いいね」だけする閲覧者）

-- users テーブルは id を明示せず採番させたいが、初期 seed が id を明示挿入した影響で
-- users_id_seq が 1 のまま進んでいない。先に現在の最大 id までシーケンスを進めておかないと
-- 自動採番が既存 id (1-4) と衝突するため、ここで補正する。
SELECT setval(
    pg_get_serial_sequence('users', 'id'),
    (SELECT MAX(id) FROM users),
    true
);

INSERT INTO users (name, password_hash) VALUES
    ('鈴木健一', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('高橋美咲', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('渡辺翔太', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('伊藤さくら', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('中村大輔', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('小林優子', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u')
ON CONFLICT (name) DO NOTHING;

-- プロフィール（アバターあり / 自己紹介あり）
INSERT INTO user_profiles (user_id, object_key, bio, is_uploaded)
SELECT u.id, v.object_key, v.bio, v.is_uploaded
FROM (VALUES
    ('鈴木健一',  'avatar/5.jpg', '神戸電子で講師をしています。Web開発とチーム開発の基礎を担当。質問は気軽にどうぞ！', TRUE),
    ('高橋美咲',  'avatar/6.jpg', 'フロントエンド志望の学生です。Vue と TypeScript を勉強中。最近は型パズルにハマっています。', TRUE),
    ('渡辺翔太',  'avatar/7.jpg', 'サーバーサイドに興味があります。Go と個人開発でアウトプット中。蔵書管理アプリを作りました。', TRUE),
    ('伊藤さくら', '',             'デザインとコーディングの両方をやりたい学生です。CSS と UI 設計を勉強しています。', FALSE),
    ('中村大輔',  '',             'プログラミングを始めて半年の初心者です。毎日コツコツ続けています。', FALSE)
) AS v(user_name, object_key, bio, is_uploaded)
JOIN users u ON u.name = v.user_name
ON CONFLICT (user_id) DO NOTHING;

-- 小林優子はあえてプロフィール行を作らない（最小状態の確認用）。
