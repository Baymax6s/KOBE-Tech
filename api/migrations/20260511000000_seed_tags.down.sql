-- 上記 seed で投入した (article, tag) の組み合わせのみを削除する。
-- ユーザーが手動で付けた本物のタグ付けを巻き込まないよう、
-- 記事は (title, user_name) で一意に識別する。
DELETE FROM article_tags at
USING articles a, users u, tags t
WHERE at.article_id = a.id
  AND at.tag_id     = t.id
  AND a.user_id     = u.id
  AND (a.title, u.name, t.name) IN (
    ('神戸大学でのハッカソン体験記',           'admin',  'ハッカソン'),
    ('神戸大学でのハッカソン体験記',           'admin',  '神戸大学'),
    ('Goで作るREST API入門',                   '田中太郎', 'Go'),
    ('Goで作るREST API入門',                   '田中太郎', 'REST API'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '山田花子', 'Vue'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '山田花子', 'フロントエンド'),
    ('PostgreSQLのマイグレーション管理',       '佐藤次郎', 'PostgreSQL'),
    ('PostgreSQLのマイグレーション管理',       '佐藤次郎', 'マイグレーション'),
    ('Dockerで開発環境を統一する',             '田中太郎', 'Docker'),
    ('Dockerで開発環境を統一する',             '田中太郎', '開発環境')
  );

-- tags の削除は、他の article_tags から参照されていないものだけに限定する。
-- これにより、ユーザーが手動で作成して別記事に紐付けたタグや、
-- 同名タグを巻き込んで消したり FK 違反で down が失敗するのを防ぐ。
DELETE FROM tags
WHERE name IN (
    'ハッカソン',
    '神戸大学',
    'Go',
    'REST API',
    'Vue',
    'フロントエンド',
    'PostgreSQL',
    'マイグレーション',
    'Docker',
    '開発環境'
) AND NOT EXISTS (
    SELECT 1 FROM article_tags at WHERE at.tag_id = tags.id
);
