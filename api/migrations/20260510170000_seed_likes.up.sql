-- likes の seed
-- 「自分がいいねした記事」「他人もいいねしている記事」が動作確認画面で
-- 同時に見えるように、各ユーザーが複数の記事へまたがっていいねしている状態を作る。
-- 自分の記事に自分でいいねは付けない。
INSERT INTO likes (article_id, user_id, created_at, updated_at)
SELECT a.id, u.id, v.created_at, v.updated_at
FROM (VALUES
    ('神戸大学でのハッカソン体験記',          '田中太郎', '2026-04-02 10:00:00+09'::timestamptz, '2026-04-02 10:00:00+09'::timestamptz),
    ('神戸大学でのハッカソン体験記',          '山田花子', '2026-04-02 11:00:00+09'::timestamptz, '2026-04-02 11:00:00+09'::timestamptz),
    ('神戸大学でのハッカソン体験記',          '佐藤次郎', '2026-04-02 12:00:00+09'::timestamptz, '2026-04-02 12:00:00+09'::timestamptz),
    ('Goで作るREST API入門',                  'admin',  '2026-04-11 09:00:00+09'::timestamptz, '2026-04-11 09:00:00+09'::timestamptz),
    ('Goで作るREST API入門',                  '山田花子', '2026-04-11 10:30:00+09'::timestamptz, '2026-04-11 10:30:00+09'::timestamptz),
    ('Goで作るREST API入門',                  '佐藤次郎', '2026-04-12 13:00:00+09'::timestamptz, '2026-04-12 13:00:00+09'::timestamptz),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'admin',  '2026-04-16 09:00:00+09'::timestamptz, '2026-04-16 09:00:00+09'::timestamptz),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '田中太郎', '2026-04-16 11:00:00+09'::timestamptz, '2026-04-16 11:00:00+09'::timestamptz),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '佐藤次郎', '2026-04-17 14:00:00+09'::timestamptz, '2026-04-17 14:00:00+09'::timestamptz),
    ('PostgreSQLのマイグレーション管理',      'admin',  '2026-04-19 09:00:00+09'::timestamptz, '2026-04-19 09:00:00+09'::timestamptz),
    ('PostgreSQLのマイグレーション管理',      '田中太郎', '2026-04-19 12:00:00+09'::timestamptz, '2026-04-19 12:00:00+09'::timestamptz),
    ('Dockerで開発環境を統一する',            'admin',  '2026-04-20 13:00:00+09'::timestamptz, '2026-04-20 13:00:00+09'::timestamptz),
    ('Dockerで開発環境を統一する',            '山田花子', '2026-04-21 09:00:00+09'::timestamptz, '2026-04-21 09:00:00+09'::timestamptz),
    ('Dockerで開発環境を統一する',            '佐藤次郎', '2026-04-21 15:00:00+09'::timestamptz, '2026-04-21 15:00:00+09'::timestamptz)
) AS v(article_title, user_name, created_at, updated_at)
JOIN articles a ON a.title = v.article_title
JOIN users u    ON u.name  = v.user_name
ON CONFLICT (article_id, user_id) DO NOTHING;
