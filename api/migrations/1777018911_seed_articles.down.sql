DELETE FROM articles
WHERE (title, user_id) IN (
    SELECT v.title, u.id
    FROM (VALUES
        ('神戸大学でのハッカソン体験記',           'admin'),
        ('Goで作るREST API入門',                   'user01'),
        ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'user02'),
        ('PostgreSQLのマイグレーション管理',         'user03'),
        ('Dockerで開発環境を統一する',              'user01')
    ) AS v(title, user_name)
    JOIN users u ON u.name = v.user_name
);
