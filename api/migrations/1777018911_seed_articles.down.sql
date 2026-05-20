DELETE FROM articles
WHERE (title, user_id) IN (
    SELECT v.title, u.id
    FROM (VALUES
        ('神戸大学でのハッカソン体験記',           'admin'),
        ('Goで作るREST API入門',                   '田中太郎'),
        ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '山田花子'),
        ('PostgreSQLのマイグレーション管理',         '佐藤次郎'),
        ('Dockerで開発環境を統一する',              '田中太郎')
    ) AS v(title, user_name)
    JOIN users u ON u.name = v.user_name
);
