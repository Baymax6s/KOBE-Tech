DELETE FROM article_tags
WHERE (article_id, tag_id) IN (
    SELECT a.id, t.id
    FROM (VALUES
        ('神戸大学でのハッカソン体験記',           'ハッカソン'),
        ('神戸大学でのハッカソン体験記',           '開発環境'),
        ('Goで作るREST API入門',                   'Go'),
        ('Goで作るREST API入門',                   'API'),
        ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'Vue'),
        ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'フロントエンド'),
        ('PostgreSQLのマイグレーション管理',         'PostgreSQL'),
        ('PostgreSQLのマイグレーション管理',         'マイグレーション'),
        ('Dockerで開発環境を統一する',              'Docker'),
        ('Dockerで開発環境を統一する',              '開発環境')
    ) AS v(article_title, tag_name)
    JOIN articles a ON a.title = v.article_title
    JOIN tags t ON t.name = v.tag_name
);

DELETE FROM tags
WHERE name IN (
    'ハッカソン',
    'Go',
    'API',
    'Vue',
    'フロントエンド',
    'PostgreSQL',
    'マイグレーション',
    'Docker',
    '開発環境'
);
