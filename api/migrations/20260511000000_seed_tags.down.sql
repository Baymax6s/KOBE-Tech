-- 上記 seed で投入した (article, tag) の組み合わせのみを削除する。
-- ユーザーが手動で付けた本物のタグ付けを巻き込まないよう、組み合わせを明示する。
DELETE FROM article_tags at
USING articles a, tags t
WHERE at.article_id = a.id
  AND at.tag_id     = t.id
  AND (a.title, t.name) IN (
    ('神戸大学でのハッカソン体験記',           'ハッカソン'),
    ('神戸大学でのハッカソン体験記',           '神戸大学'),
    ('Goで作るREST API入門',                   'Go'),
    ('Goで作るREST API入門',                   'REST API'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'Vue'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'フロントエンド'),
    ('PostgreSQLのマイグレーション管理',       'PostgreSQL'),
    ('PostgreSQLのマイグレーション管理',       'マイグレーション'),
    ('Dockerで開発環境を統一する',             'Docker'),
    ('Dockerで開発環境を統一する',             '開発環境')
  );

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
);
