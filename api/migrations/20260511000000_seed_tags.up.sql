-- tags の seed
-- 動作確認画面で「タグ付き記事一覧」「タグでの絞り込み」が機能するように、
-- 既存の seed 記事（1777018911_seed_articles.up.sql）と紐付けられるタグを投入する。
INSERT INTO tags (name, created_at, updated_at)
SELECT v.name, v.created_at, v.updated_at
FROM (VALUES
    ('ハッカソン',   '2026-04-01 09:00:00+09'::timestamptz, '2026-04-01 09:00:00+09'::timestamptz),
    ('神戸大学',     '2026-04-01 09:00:00+09'::timestamptz, '2026-04-01 09:00:00+09'::timestamptz),
    ('Go',           '2026-04-10 12:30:00+09'::timestamptz, '2026-04-10 12:30:00+09'::timestamptz),
    ('REST API',     '2026-04-10 12:30:00+09'::timestamptz, '2026-04-10 12:30:00+09'::timestamptz),
    ('Vue',          '2026-04-15 15:00:00+09'::timestamptz, '2026-04-15 15:00:00+09'::timestamptz),
    ('フロントエンド', '2026-04-15 15:00:00+09'::timestamptz, '2026-04-15 15:00:00+09'::timestamptz),
    ('PostgreSQL',   '2026-04-18 10:00:00+09'::timestamptz, '2026-04-18 10:00:00+09'::timestamptz),
    ('マイグレーション', '2026-04-18 10:00:00+09'::timestamptz, '2026-04-18 10:00:00+09'::timestamptz),
    ('Docker',       '2026-04-20 11:00:00+09'::timestamptz, '2026-04-20 11:00:00+09'::timestamptz),
    ('開発環境',     '2026-04-20 11:00:00+09'::timestamptz, '2026-04-20 11:00:00+09'::timestamptz)
) AS v(name, created_at, updated_at)
ON CONFLICT (name) DO NOTHING;

-- article_tags の seed
-- 既存 seed 記事に対し、内容に合致するタグを 2 件ずつ紐付ける。
-- 記事は (title, user_id) で一意に識別する（同名タイトルの別ユーザー記事に誤って紐付かないように user_name も突き合わせる）。
INSERT INTO article_tags (article_id, tag_id, created_at)
SELECT a.id, t.id, v.created_at
FROM (VALUES
    ('神戸大学でのハッカソン体験記',           'admin',  'ハッカソン',     '2026-04-01 09:00:00+09'::timestamptz),
    ('神戸大学でのハッカソン体験記',           'admin',  '神戸大学',       '2026-04-01 09:00:00+09'::timestamptz),
    ('Goで作るREST API入門',                   'user01', 'Go',             '2026-04-10 12:30:00+09'::timestamptz),
    ('Goで作るREST API入門',                   'user01', 'REST API',       '2026-04-10 12:30:00+09'::timestamptz),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'user02', 'Vue',            '2026-04-15 15:00:00+09'::timestamptz),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'user02', 'フロントエンド', '2026-04-15 15:00:00+09'::timestamptz),
    ('PostgreSQLのマイグレーション管理',       'user03', 'PostgreSQL',     '2026-04-18 10:00:00+09'::timestamptz),
    ('PostgreSQLのマイグレーション管理',       'user03', 'マイグレーション', '2026-04-18 10:00:00+09'::timestamptz),
    ('Dockerで開発環境を統一する',             'user01', 'Docker',         '2026-04-20 11:00:00+09'::timestamptz),
    ('Dockerで開発環境を統一する',             'user01', '開発環境',       '2026-04-20 11:00:00+09'::timestamptz)
) AS v(article_title, user_name, tag_name, created_at)
JOIN users u    ON u.name  = v.user_name
JOIN articles a ON a.title = v.article_title AND a.user_id = u.id
JOIN tags t     ON t.name  = v.tag_name
ON CONFLICT (article_id, tag_id) DO NOTHING;
